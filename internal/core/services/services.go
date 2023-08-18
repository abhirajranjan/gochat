package services

import (
	"gochat/internal/core/ports"
)

type service struct {
	repo ports.Repositories
}

var _ ports.Service = (*service)(nil)

func NewService(repo ports.Repositories) *service {
	return &service{
		repo: repo,
	}
}
