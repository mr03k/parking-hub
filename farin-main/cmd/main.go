package main

import (
	"context"
	"farin/cmd/seeder"
	"farin/infrastructure/godotenv"
	gormdb "farin/infrastructure/gorm"
	"farin/infrastructure/opentelemetry"
	"farin/infrastructure/rabbit"
	"flag"
	"fmt"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	uploader "github.com/mahdimehrabi/uploader/minio"
	"log"
	"log/slog"
	"os"
	"time"
)

// @securityDefinitions.apikey ApiKeyAuth

func main() {
	seed := flag.Bool("seed", false, "Seed the database with initial data")
	fakeData := flag.Bool("fakeData", false, "Seed the database with initial data")
	logger := initSlogLogger()

	flag.Parse()
	if *seed {
		if err := seeder.Seed(*fakeData); err != nil {
			logger.Error("Failed to seed the database with initial data", "error", err)
			os.Exit(0) //success because we use this section on devops tools as well
		}
		return
	}

	env := godotenv.NewEnv()
	env.Load()

	fmt.Println(env)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	secure := env.Environment != "development"
	minio := uploader.NewMinio(uploader.WithConfig(
		&uploader.Config{
			Secure:           secure,
			MinioHost:        env.MinioHost,
			MinioSecret:      env.MinioSecret,
			MinioAccessToken: env.MinioAccessToken,
			Buckets:          []string{env.MinioProfilePictureBucket},
		},
	))
	if err := minio.Setup(ctx); err != nil {
		log.Fatalf("failed to setup minio:%s", err)
	}

	ot := opentelemetry.NewOpenTelemetry(env)
	if err := ot.Setup(ctx); err != nil {
		log.Fatalf("failed to setup opentelemetry:%s", err)
	}
	defer ot.Shutdown(ctx)
	logger = slog.New(slogmulti.Fanout(
		otelslog.NewHandler("ot-slog", otelslog.WithLoggerProvider(ot.Logger)),
		slog.NewTextHandler(os.Stdout, nil),
	))

	rbt := rabbit.NewRabbit(env, logger)
	gorm := gormdb.NewGORMDB(env, logger)
	if err := gorm.Setup(ctx); err != nil {
		log.Fatalf("failed to setup gorm:%s", err)
	}

	boot, err := wireApp(logger, rbt, minio, gorm, ot)
	if err != nil {
		log.Fatalf("failed to setup app:%s", err)
	}
	boot.Boot()
}

func initSlogLogger() *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions))
	slog.SetDefault(logger)

	return logger
}
