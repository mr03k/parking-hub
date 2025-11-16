package gormdb

import (
	"context"
	"errors"
	"fmt"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMDB struct {
	DB  *gorm.DB
	env *godotenv.Env
}

func NewGORMDB(env *godotenv.Env) *GORMDB {
	return &GORMDB{
		env: env,
	}
}

func (g *GORMDB) Setup(ctx context.Context) error {
	// Close existing connection if present
	if g.DB != nil {
		sqlDB, err := g.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	// Configure GORM connection
	dsn := g.env.DatabaseHost

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB from gorm: %w", err)
	}
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	g.DB = db
	return nil
}

// HealthCheck verifies the PostgreSQL database connection
func (g *GORMDB) HealthCheck(ctx context.Context) error {
	if g.DB == nil {
		return errors.New("gorm database is not initialized")
	}

	sqlDB, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB from gorm: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ping database to verify connection
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Close releases all database connections
func (g *GORMDB) Close() {
	if g.DB != nil {
		sqlDB, err := g.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
