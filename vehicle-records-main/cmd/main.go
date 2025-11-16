package main

import (
	"context"
	"flag"
	"git.abanppc.com/farin-project/vehicle-records/cmd/seeder"
	"git.abanppc.com/farin-project/vehicle-records/domain/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	gormdb "git.abanppc.com/farin-project/vehicle-records/infrastructure/gorm"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/rabbit"
	uploader "github.com/mahdimehrabi/uploader/minio"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"log"
	"log/slog"
	"os"
	"time"
)

// @securityDefinitions.apikey ApiKeyAuth

func main() {
	seed := flag.Bool("seed", false, "Seed the database with initial data")
	retry := flag.Bool("retry", false, "retry from database")
	fromTime := flag.Int64("from", 0, "retry from time")
	limit := flag.Int("limit", 100, "limit retry count")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secure := env.Environment != "development"
	minio := uploader.NewMinio(uploader.WithConfig(
		&uploader.Config{
			Secure:           secure,
			MinioHost:        env.MinioHost,
			MinioSecret:      env.MinioSecret,
			MinioAccessToken: env.MinioAccessToken,
			Buckets:          []string{env.MinioVehicleRecordsBucket},
			DisableMultiPart: true,
		},
	))
	if err := minio.Setup(ctx); err != nil {
		log.Fatalf("failed to setup minio:%s", err)
	}

	ot := opentelemetry.NewOpenTelemetry(env)
	if err := ot.Setup(ctx); err != nil {
		log.Fatalf("failed to setup opentelemetry:%s", err)
	}
	logger = slog.New(slogmulti.Fanout(
		otelslog.NewHandler("ot-slog", otelslog.WithLoggerProvider(ot.Logger)),
		slog.NewTextHandler(os.Stdout, nil),
	))

	rbt := rabbit.NewRabbit(env, logger)
	gorm := gormdb.NewGORMDB(env)
	if err := gorm.Setup(ctx); err != nil {
		log.Fatalf("failed to setup gorm:%s", err)
	}

	boot, err := wireApp(logger, rbt, minio, gorm, ot)
	if err != nil {
		log.Fatalf("failed to setup app:%s", err)
	}
	boot.Boot(*retry, *fromTime, *limit)
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
