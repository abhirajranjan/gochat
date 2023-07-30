package handlers

import (
	"gochat/internal/core/domain"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func (h *handler) HandleGoogleAuth(ctx *gin.Context) {
	var (
		cred         Credential
		claims       domain.GoogleJWTModel
		loginRequest domain.LoginRequest
		sessionjwt   domain.SessionJwtModel
	)

	if err := ctx.Bind(&cred); err != nil {
		h.logger.Debugf("HandleGoogleAuth: error binding context, %w", err)
		setUnauthorised(ctx)
		return
	}

	_, _, err := h.jwtParser.ParseUnverified(cred.Token(), &claims)
	if err != nil {
		h.logger.Debugf("HandleGoogleAuth: error parsing token %w", err)
		setInvalidToken(ctx)
		return
	}

	loginRequest = domain.LoginRequest{
		Email:       claims.Email,
		Name:        claims.Name,
		Picture:     claims.Picture,
		Given_name:  claims.Given_name,
		Family_name: claims.Family_name,
	}

	user, err := h.service.LoginRequest(loginRequest)
	if err != nil {
		h.logger.Errorf("HandleGoogleAuth: %w", err)
		setInternalServerError(ctx)
		return
	}

	sessionjwt = domain.SessionJwtModel{
		JwtModel: domain.JwtModel{
			Iss: "dailydsa/auth/v1",
			Aud: []string{strconv.FormatInt(user.UserID, 10)},
			Nbf: jwt.NewNumericDate(time.Now()),
			Exp: jwt.NewNumericDate(time.Now().Add(h.config.Expiry)),
		},

		UserID:     user.UserID,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Email:      user.Email,
		Picture:    user.Picture,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionjwt)
	jwt, err := token.SignedString([]byte(h.config.Key))
	if err != nil {
		h.logger.Errorf("HandleGoogleAuth: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
	h.logger.Debugf("response jwt: %s", jwt)
}

// injects user information from jwt to context
func (h *handler) AuthMiddleware(ctx *gin.Context) {
	var session domain.SessionJwtModel

	jwtoken, err := extractTokenFromCtx(ctx)
	if err != nil {
		h.logger.Debugf("authMiddleware: %w", err)
		setUnauthorised(ctx)
		return
	}

	if _, err := h.jwtParser.ParseWithClaims(jwtoken, &session, func(j *jwt.Token) (interface{}, error) {
		// get hashing key according to algorithm
		if j.Method.Alg() == jwt.SigningMethodHS256.Name {
			return h.config.Key, nil
		}
		// no algorithm matched
		return nil, errors.New("wrong algorithm")

	}); err != nil {
		h.logger.Debugf("authMiddleware: %w", err)
		setInvalidToken(ctx)
		return
	}

	ctx.Set("userId", session.UserID)
	ctx.Next()
}

func extractTokenFromCtx(ctx *gin.Context) (jwt string, err error) {
	// header based token extraction
	token := ctx.Request.Header.Get("Authorisation")
	b, jwt, ok := strings.Cut(token, "Bearer ")
	if !ok {
		return "", errors.New("No authorisation header found")
	}

	if b != "" {
		return "", errors.New("Bearer token mismatched")
	}

	return jwt, nil
}
