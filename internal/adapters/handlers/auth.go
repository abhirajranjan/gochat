package handlers

import (
	"gochat/internal/core/domain"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// injects user information from jwt to context
func (h *handler) AuthMiddleware(ctx *gin.Context) {
	// TODO: token expired
	var session domain.SessionJwtModel

	jwtoken, err := extractTokenFromCtx(ctx)
	if err != nil {
		h.Debugf("extractTokenFromCtx: %w", err)
		setForbidden(ctx)
		return
	}

	if _, err := h.jwtParser.ParseWithClaims(jwtoken, &session, func(j *jwt.Token) (interface{}, error) {
		return h.getEncryptionKey(j.Method)
	}); err != nil {
		h.Debugf("jwtParser.ParseWithClaims: %w", err)
		setInvalidToken(ctx)
		return
	}

	ctx.Set(NAMETAGKEY, session.Sub)
	ctx.Next()
}

func (h *handler) parseUnverified(String string, claims jwt.Claims) (token *jwt.Token, parts []string, err error) {
	return h.jwtParser.ParseUnverified(String, claims)
}

func (h *handler) generateSessionJwt(user *domain.User) (jwttoken string, err error) {
	sessionjwt := domain.SessionJwtModel{
		JwtModel: domain.JwtModel{
			Iss: JWT_ISSUER,
			Aud: []string{user.NameTag},
			Sub: user.ID,
			Nbf: jwt.NewNumericDate(time.Now()),
			Exp: jwt.NewNumericDate(time.Now().Add(h.config.Expiry)),
		},

		NameTag:    user.NameTag,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Email:      user.Email,
		Picture:    user.Picture,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionjwt)
	key, err := h.getEncryptionKey(jwt.SigningMethodHS256)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

// get hashing key according to algorithm
func (h *handler) getEncryptionKey(method jwt.SigningMethod) (interface{}, error) {
	if method == jwt.SigningMethodHS256 {
		return []byte(h.config.Key), nil
	}
	// no algorithm matched
	return nil, errors.New("wrong algorithm")
}

func extractTokenFromCtx(ctx *gin.Context) (id string, err error) {
	// header based token extraction
	token := ctx.Request.Header.Get("Authorization")
	b, id, ok := strings.Cut(token, "Bearer ")
	if !ok {
		return "", errors.New("No authorization header found")
	}

	if b != "" {
		return "", errors.New("Bearer token mismatched")
	}

	return id, nil
}
