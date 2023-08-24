package services

import (
	"fmt"
	"gochat/internal/core/domain"
)

var (
	ErrChannelNotFound = domain.NewErrDomain("invalid channel")
	ErrUserNotFound    = domain.NewErrDomain("invalid user")
)

func ErrCannotBeEmpty(cantbeempty string) error {
	return domain.NewErrDomain(fmt.Sprintf("%s cannot be empty", cantbeempty))
}
