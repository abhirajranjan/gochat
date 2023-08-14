package mysql

import (
	"gochat/config"
	"gochat/internal/adapters/repositories/sql"
	"gochat/internal/core/ports"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySqlRepository(cfg config.SqlConfig) (ports.Repositories, error) {
	return sql.NewSqlRepository(cfg, mysql.Open(cfg.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
