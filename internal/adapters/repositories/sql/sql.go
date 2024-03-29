package sql

import (
	"context"
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
	return &sqlRepo{
		config: cfg,
		conn:   conn,
	}, nil
}

func (r *sqlRepo) getContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, r.config.SqlTimeout)
}
