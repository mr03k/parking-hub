package gormdb

import (
	"context"
	"errors"
	"farin/infrastructure/godotenv"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger implements GORM's logger.Interface using slog
type GormLogger struct {
	slogger *slog.Logger
	config  logger.Config
}

// NewGormLogger creates a new GORM logger that uses slog
func NewGormLogger(slogger *slog.Logger) *GormLogger {
	return &GormLogger{
		slogger: slogger,
		config: logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	}
}

// LogMode implements GORM's logger interface
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level
	return &newLogger
}

// Info implements GORM's logger interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.slogger.InfoContext(ctx, fmt.Sprintf(msg, data...))
}

// Warn implements GORM's logger interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.slogger.WarnContext(ctx, fmt.Sprintf(msg, data...))
}

// Error implements GORM's logger interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.slogger.ErrorContext(ctx, fmt.Sprintf(msg, data...))
}

// Trace implements GORM's logger interface
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []any{
		slog.String("elapsed", elapsed.String()),
		slog.String("sql", sql),
		slog.Int64("rows", rows),
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		attrs = append(attrs, slog.Any("error", err))
		l.slogger.ErrorContext(ctx, "gorm-trace", attrs...)
		return
	}

	if elapsed > l.config.SlowThreshold {
		l.slogger.WarnContext(ctx, "gorm-slow-query", attrs...)
		return
	}

	l.slogger.DebugContext(ctx, "gorm-trace", attrs...)
}

type GORMDB struct {
	DB     *gorm.DB
	env    *godotenv.Env
	logger *slog.Logger
}

func NewGORMDB(env *godotenv.Env, logger *slog.Logger) *GORMDB {
	return &GORMDB{
		env:    env,
		logger: logger,
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

	dsn := g.env.DatabaseHost

	gormLogger := NewGormLogger(g.logger)
	gormLogger.LogMode(logger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
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

func (g *GORMDB) Close() {
	if g.DB != nil {
		sqlDB, err := g.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
