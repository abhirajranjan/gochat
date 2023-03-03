package AuthMiddleware

import jwt "github.com/appleboy/gin-jwt/v2"

var (
	ErrMissingLoginValues   = jwt.ErrMissingLoginValues
	ErrFailedAuthentication = jwt.ErrFailedAuthentication
)
