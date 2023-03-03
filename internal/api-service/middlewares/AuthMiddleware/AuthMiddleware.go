package AuthMiddleware

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/abhirajranjan/gochat/pkg/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type jwtAuth struct {
	*jwt.GinJWTMiddleware
	handler IHandler
	Logger  logger.ILogger
	Cfg     AuthConf
}

func NewJWTMiddleware(cfg AuthConf, logger logger.ILogger, methodhandler IHandler) (*jwtAuth, error) {
	jwtauth := &jwtAuth{
		Cfg:     cfg,
		Logger:  logger,
		handler: methodhandler,
	}

	jwtmw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         cfg.Realm,
		Key:           cfg.Key,
		Timeout:       cfg.TimeoutDuration,
		MaxRefresh:    cfg.MaxRefresh,
		IdentityKey:   cfg.IdentityKey,
		TokenHeadName: cfg.TokenHeadName,
		TokenLookup:   cfg.TokenLookup,

		HTTPStatusMessageFunc: jwtauth.HTTPStatusMessageFunc,
		Authorizator:          jwtauth.authorizator,
		PayloadFunc:           jwtauth.payloadFunc,
		LoginResponse:         jwtauth.loginResponse,
		IdentityHandler:       jwtauth.identityHandler,
		Authenticator:         jwtauth.authenticator,
		RefreshResponse:       jwtauth.refreshResponse,
		LogoutResponse:        jwtauth.LogoutResponse,
		Unauthorized:          jwtauth.unauthorized,
		TimeFunc:              jwtauth.timefunc,
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt.New")
	}
	jwtauth.GinJWTMiddleware = jwtmw
	return jwtauth, nil
}

func (j *jwtAuth) authenticator(c *gin.Context) (interface{}, error) {
	loginreq, err := j.handler.GenerateLoginRequest(c)
	if err != nil {
		return nil, err
	}

	loginres, err := j.handler.HandleLoginRequest(loginreq)
	if err != nil {
		return nil, err
	}

	return loginres, nil
}

func (j *jwtAuth) HTTPStatusMessageFunc(e error, c *gin.Context) string {
	if b, ok := e.(serviceErrors.IErr); ok {
		bytes, err := b.To_json()
		if err != nil {
			j.Logger.Error(errors.Wrap(err, "jwtauth.httpSatusMessageFunc"))
			return ""
		}
		return string(bytes)
	}
	if errors.Is(e, serviceErrors.ErrInternalServer) {
		return "internal server error"
	}
	j.Logger.Error(errors.Wrap(e, "jwtauth.httpSatusMessageFunc"))
	return e.Error()
}

func (j *jwtAuth) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(ILoginResponse); ok {
		return j.handler.GeneratePayloadData(v)
	}
	j.Logger.Error(serviceErrors.NewStandardErr("jwtAuth.payloadFunc", "cannot decode unmarshal into ILoginResponse", data))
	return jwt.MapClaims{}
}

func (j *jwtAuth) loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	j.Logger.Debug("login successfull, token:%s expiry:%v", token, expire)
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func (j *jwtAuth) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return j.handler.ExtractPayloadData(claims)
}

func (j *jwtAuth) authorizator(data interface{}, c *gin.Context) bool {

	return false
}

func (j *jwtAuth) refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	j.Logger.Infof(`refresh token "%s" passed [%d], new expire on: %v`, token, code, expire)
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func (j *jwtAuth) LogoutResponse(c *gin.Context, code int) {

}

func (j *jwtAuth) unauthorized(c *gin.Context, code int, message string) {
	j.Logger.Warnf(`unauthorized jwt token "%s" passed [%d]: %s`, jwt.GetToken(c), code, message)
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (j *jwtAuth) timefunc() time.Time {
	return time.Now()
}
