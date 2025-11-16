package main

import (
	"context"
	"flag"
	"git.abanppc.com/farin-project/state/cmd/seeder"
	"git.abanppc.com/farin-project/state/infrastructure/godotenv"
	gormdb "git.abanppc.com/farin-project/state/infrastructure/gorm"
	"git.abanppc.com/farin-project/state/infrastructure/rabbit"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rbt := rabbit.NewRabbit(env, logger)
	gorm := gormdb.NewGORMDB(env)
	if err := gorm.Setup(ctx); err != nil {
		log.Fatalf("failed to setup gorm:%s", err)
	}

	boot, err := wireApp(logger, rbt, gorm)
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
