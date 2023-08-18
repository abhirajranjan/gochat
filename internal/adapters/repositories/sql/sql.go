package sql

import (
	"gochat/config"
	"gochat/internal/core/ports"

	"gorm.io/gorm"
)

type sqlRepo struct {
	config config.SqlConfig
	conn   *gorm.DB
}

var _ ports.Repositories = (*sqlRepo)(nil)

func NewSqlRepository(cfg config.SqlConfig, dialector gorm.Dialector, opts ...gorm.Option) (*sqlRepo, error) {
	conn, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}

	conn.AutoMigrate(&User{}, &UserChannels{}, &Messages{}, &Channel{})
	conn.Migrator().CreateTable(&User{}, &UserChannels{}, &Messages{}, &Channel{})

	return &sqlRepo{
		config: cfg,
		conn:   conn,
	}, nil
}
