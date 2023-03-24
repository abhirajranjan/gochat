package JwtAuthMiddleware

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/abhirajranjan/gochat/pkg/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var jwterr = []error{jwt.ErrForbidden,
	jwt.ErrFailedTokenCreation, jwt.ErrExpiredToken, jwt.ErrEmptyAuthHeader,
	jwt.ErrMissingExpField, jwt.ErrWrongFormatOfExp, jwt.ErrInvalidAuthHeader,
	jwt.ErrEmptyQueryToken, jwt.ErrEmptyCookieToken, jwt.ErrEmptyParamToken}

type jwtAuth struct {
	jwt     *jwt.GinJWTMiddleware
	handler model.IHandler
	Logger  logger.ILogger
	Cfg     *AuthConf
}

func NewJWTMiddleware(cfg *AuthConf, logger logger.ILogger, methodhandler model.IHandler) (*jwtAuth, error) {
	jwtauth := &jwtAuth{
		Cfg:     cfg,
		Logger:  logger,
		handler: methodhandler,
	}

	jwtmw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         cfg.Realm,
		Key:           []byte(cfg.Key),
		Timeout:       cfg.TimeoutDuration,
		MaxRefresh:    cfg.MaxRefresh,
		IdentityKey:   cfg.IdentityKey,
		TokenHeadName: cfg.TokenHeadName,
		TokenLookup:   cfg.TokenLookup,

		// login
		Authenticator: jwtauth.authenticator,
		PayloadFunc:   jwtauth.payloadFunc,
		LoginResponse: jwtauth.loginResponse,

		IdentityHandler: jwtauth.identityHandler,
		Authorizator:    jwtauth.authorizator,

		LogoutResponse:        jwtauth.logoutResponse,
		RefreshResponse:       jwtauth.refreshResponse,
		HTTPStatusMessageFunc: jwtauth.httpStatusMessageFunc,
		Unauthorized:          jwtauth.unauthorized,
		TimeFunc:              jwtauth.timefunc,
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt.New")
	}

	jwtauth.jwt = jwtmw
	return jwtauth, nil
}

//* login

// login middleware for login action
func (j *jwtAuth) Login() gin.HandlerFunc {
	return j.jwt.LoginHandler
}

// extract login credentials from gin.Context and returns user data that needs to be embeded into jwt.
//
// data return by authenticator is passed to payload function to convert data into mapClaims
// if error is returned, unauthorised is called
func (j *jwtAuth) authenticator(c *gin.Context) (interface{}, error) {
	loginres, err := j.handler.HandleLoginRequest(c)
	if err != nil {
		return nil, err
	}

	return loginres, nil
}

// take whatever was returned from Authenticator and convert it into MapClaims (i.e. map[string]interface{}).
// set all the values that are needed to be embeded in the jwt token
func (j *jwtAuth) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(model.ILoginResponse); ok {
		return j.handler.GeneratePayloadData(v)
	}
	j.Logger.Warn(serviceErrors.NewStandardErr("jwtAuth.payloadFunc", "cannot decode unmarshal into ILoginResponse", data))
	return jwt.MapClaims{}
}

// set user response when login response was successful.
func (j *jwtAuth) loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	j.Logger.Debugf("login successful, token: %s expiry: %v", token, expire)
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

//* check authentication

// provide access to middleware function
func (j *jwtAuth) CheckIfValidAuth() gin.HandlerFunc {
	return j.jwt.MiddlewareFunc()
}

// fetch the user identity from claims embedded within the jwt token, and pass this identity value to Authorizator
func (j *jwtAuth) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	var payload model.IPayloadData = j.handler.ExtractPayloadData(claims)
	return payload
}

// check if the user is authorized to access endpoint
//
// return true if the user is authorized to continue through with the request, or false if they are not authorized
// in case of failure, unauthorized is called
func (j *jwtAuth) authorizator(data interface{}, c *gin.Context) bool {
	payload := data.(model.IPayloadData)
	return j.handler.VerifyUser(payload)
}

//* logout

// logout middleware for logout action
//
// calls logoutResponse to set response if successful or not
func (j *jwtAuth) Logout() gin.HandlerFunc {
	return j.jwt.LogoutHandler
}

// called when logout action was called and returns if successful or not
func (j *jwtAuth) logoutResponse(c *gin.Context, code int) {
	claims := jwt.ExtractClaims(c)
	status := j.handler.LogoutUser(claims)
	c.JSON(status, gin.H{
		"code": status,
	})
}

//* refresh token

// refresh token middleware function
//
// if refresh token in passed then refreshResponse is called else unauthorised is called
func (j *jwtAuth) RefreshToken() gin.HandlerFunc {
	return j.jwt.RefreshHandler
}

// called when refresh token action is passed
func (j *jwtAuth) refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	j.Logger.Infof(`refresh token "%s" passed [%d], new expire on: %v`, token, code, expire)
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

//* other functions

// function to return message to write into json response if error was raised
func (j *jwtAuth) httpStatusMessageFunc(e error, c *gin.Context) string {
	// check if custom error of type serviceErrors IErr
	if b, ok := e.(serviceErrors.IErr); ok {
		// service errors contains To_json that convert respective error call to error json
		bytes, err := b.To_json()
		if err != nil {
			j.Logger.Warn(errors.Wrap(err, "jwtauth.httpStatusMessageFunc"))
			return "internal server error"
		}
		return string(bytes)
	}

	if errors.Is(e, serviceErrors.ErrInternalServer) {
		return "internal server error"
	}

	// check if error is jwt module error
	for _, err := range jwterr {
		if errors.Is(e, err) {
			return err.Error()
		}
	}

	j.Logger.Warn(errors.Wrap(e, "jwtauth.httpStatusMessageFunc"))
	return "internal server error"
}

// called when failures with logging in, bad tokens, or lacking privileges
func (j *jwtAuth) unauthorized(c *gin.Context, code int, message string) {
	// logs the error
	j.Logger.Debugf(`jwt token "%s" failed [%d]: %s`, jwt.GetToken(c), code, message)
	// return response in form of code and message
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// custom time function to check for time of expiry
func (j *jwtAuth) timefunc() time.Time {
	return time.Now()
}
