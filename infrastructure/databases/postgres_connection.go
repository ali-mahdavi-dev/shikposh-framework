package databases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ali-mahdavi-dev/framework/infrastructure/logging"
)

type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxIdleTime  int
	MaxOpenConns int
	MaxIdleConns int
}

func New(cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(cfg.DBType) {
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite3":
		_ = os.MkdirAll(filepath.Dir(cfg.DSN), os.ModePerm)
		dialector = sqlite.Open(cfg.DSN)
	default:
		err := fmt.Errorf("unsupported database type: %s", cfg.DBType)
		logging.Error("Unsupported database type").
			WithString("db_type", cfg.DBType).
			WithError(err).
			Log()
		return nil, err
	}

	ormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Discard,
	}

	if cfg.Debug {
		ormCfg.Logger = logger.Default
	}

	db, err := gorm.Open(dialector, ormCfg)
	if err != nil {
		logging.Error("Failed to open database connection").
			WithString("db_type", cfg.DBType).
			WithError(err).
			Log()
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		logging.Error("Failed to get underlying SQL database").
			WithError(err).
			Log()
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		logging.Error("Database ping failed").
			WithError(err).
			Log()
		return nil, err
	}

	// Log database connection pool configuration (only once at startup, not performance critical)
	logging.Info("Database connection pool configured").
		WithString("db_type", cfg.DBType).
		WithInt("max_open_conns", cfg.MaxOpenConns).
		WithInt("max_idle_conns", cfg.MaxIdleConns).
		WithInt("max_lifetime_sec", cfg.MaxLifetime).
		WithInt("max_idle_time_sec", cfg.MaxIdleTime).
		Log()

	return db, nil
}
