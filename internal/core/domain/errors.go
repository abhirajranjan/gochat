package domain

var (
// ErrDomain = errors.New("domain error")
)

type ErrDomain struct {
	reason string
}

func (e ErrDomain) Reason() string {
	return e.reason
}

func (e ErrDomain) Error() string {
	return e.reason
}

func NewErrDomain(cause string) error {
	return ErrDomain{reason: cause}
}
