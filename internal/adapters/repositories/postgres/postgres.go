package postgres

import (
	"gochat/config"
	"gochat/internal/adapters/repositories/sql"
	"gochat/internal/core/ports"

	"gorm.io/driver/postgres"
)

func NewPostgresRepository(cfg config.SqlConfig) (ports.Repositories, error) {
	return sql.NewSqlRepository(cfg, postgres.Open(cfg.DSN))
}
