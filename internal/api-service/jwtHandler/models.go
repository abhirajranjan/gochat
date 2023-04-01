package jwtHandler

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginRequest) GetUsername() string {
	return l.Username
}

func (l *LoginRequest) SetUsername(username string) {
	l.Username = username
}

func (l *LoginRequest) GetPassword() string {
	return l.Password
}

func (l *LoginRequest) SetPassword(password string) {
	l.Password = password
}
