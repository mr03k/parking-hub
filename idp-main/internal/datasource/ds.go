package datasource

import (
	"database/sql"
	"log/slog"

	"application/config"

	_ "github.com/lib/pq"
)

type Datasource struct {
	logger *slog.Logger
	DBpsql *sql.DB
	cfg    config.Config
}

// New creates a new datasource
func NewDataSource(cfg config.Config, logger *slog.Logger) *Datasource {
	ds := &Datasource{
		cfg:    cfg,
		logger: logger.With("layer", "Datasource"),
	}
	if err := ds.initPostgres(); err != nil {
		ds.logger.Error("initPostgres", "error", err)
		panic(err)
	}
	return ds
}

type postgressConfig struct {
	DSN     string `koanf:"dsn"`
	Enabled bool   `koanf:"enabled"`
}

// init postgress

func (ds *Datasource) initPostgres() error {
	cfg := new(postgressConfig)
	ds.cfg.Unmarshal("datasource.postgres", cfg)
	ds.logger.Debug("initPostgres", "config", cfg)

	db, err := sql.Open("postgres", cfg.DSN)
	// db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	ds.logger.Info("initPostgres", "status", "success")
	ds.DBpsql = db
	return nil
}
