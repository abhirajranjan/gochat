package handlers

type Credential struct {
	Token string `json:"token" binding:"required"`
}

func (c Credential) GetToken() string {
	return c.Token
}
