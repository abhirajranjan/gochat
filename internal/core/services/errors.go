package services

import (
	"gochat/internal/core/ports"

	"github.com/pkg/errors"
)

var (
	ErrChannelNotFound  = errors.Wrap(ports.ErrDomain, "invalid channel")
	ErrUserNotInChannel = errors.Wrap(ports.ErrDomain, "user not in channel")

	ErrUserNotFound = errors.Wrap(ports.ErrDomain, "invalid user")
)

func ErrCannotBeEmpty(cantbeempty string) error {
	return errors.Wrapf(ports.ErrDomain, "%s cannot be empty", cantbeempty)
}
