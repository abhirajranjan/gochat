package grpcHandler

type ILoginRequest interface {
	GetUsername() string
	SetUsername(string)

	GetPassword() string
	SetPassword(string)
}

type LoginRequest struct {
	username string
	password string
}

func (l *LoginRequest) GetUsername() string {
	return l.username
}

func (l *LoginRequest) SetUsername(username string) {
	l.username = username
}

func (l *LoginRequest) GetPassword() string {
	return l.password
}

func (l *LoginRequest) SetPassword(password string) {
	l.password = password
}

//////////////////////////////////////////////////////////////////////////////////////////////

type ILoginResponse interface {
	GetUserID() string
	GetUserRoles() []int64
	GetErr() error
	GetErrCode() int64
}

type IPayloadData interface {
	GetUserID() string
	GetUserRoles() []int64
}

type responseStatus struct {
	err     error
	errCode int64
}

type loginResponse struct {
	UserID    string  `payload:"userID"`
	UserRoles []int64 `payload:"userRoles"`
	responseStatus
}

func (l *loginResponse) GetUserID() string {
	return l.UserID
}

func (l *loginResponse) GetUserRoles() []int64 {
	return l.UserRoles
}

func (l *loginResponse) GetErr() error {
	return l.err
}

func (l *loginResponse) GetErrCode() int64 {
	return l.errCode
}

func NewLoginResponse(userID string, userRoles []int64, errCode int64, err error) ILoginResponse {
	return &loginResponse{
		UserID:    userID,
		UserRoles: userRoles,

		responseStatus: responseStatus{
			err:     err,
			errCode: errCode,
		},
	}
}
