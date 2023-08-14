package handlers

type Credential struct {
	Token string `json:"token"`
}

func (c Credential) GetToken() string {
	return c.Token
}
