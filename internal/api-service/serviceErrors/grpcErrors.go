package serviceErrors

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
)
