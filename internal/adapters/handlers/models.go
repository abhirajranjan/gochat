package handlers

type Credential struct {
	token string
}

func (c Credential) Token() string {
	return c.token
}
