package services

import (
	"fmt"
	"gochat/internal/core/domain"
)

var (
	ErrChannelNotFound  = domain.NewErrDomain("invalid channel")
	ErrUserNotInChannel = domain.NewErrDomain("user not in channel")
	ErrUserNotFound     = domain.NewErrDomain("invalid user")
)

func ErrCannotBeEmpty(cantbeempty string) error {
	return domain.NewErrDomain(fmt.Sprintf("%s cannot be empty", cantbeempty))
}
