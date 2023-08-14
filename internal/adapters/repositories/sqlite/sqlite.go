package sqlite

import (
	"gochat/config"
	"gochat/internal/adapters/repositories/sql"
	"gochat/internal/core/ports"

	"gorm.io/driver/sqlite"
)

func NewSqliteRepository(cfg config.SqlConfig) (ports.Repositories, error) {
	return sql.NewSqlRepository(cfg, sqlite.Open(cfg.DSN))
}