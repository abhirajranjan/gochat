package handlers

import (
	"gochat/internal/core/domain"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// injects user information from jwt to context
func (h *handler) AuthMiddleware(ctx *gin.Context) {
	var session domain.SessionJwtModel

	jwtoken, err := extractTokenFromCtx(ctx)
	if err != nil {
		h.Debugf("extractTokenFromCtx: %w", err)
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
		h.Debugf("jwtParser.ParseWithClaims: %w", err)
		setInvalidToken(ctx)
		return
	}

	ctx.Set(NAMETAGKEY, session.Sub)
	ctx.Next()
}

func extractTokenFromCtx(ctx *gin.Context) (id string, err error) {
	// header based token extraction
	token := ctx.Request.Header.Get("Authorisation")
	b, id, ok := strings.Cut(token, "Bearer ")
	if !ok {
		return "", errors.New("No authorisation header found")
	}

	if b != "" {
		return "", errors.New("Bearer token mismatched")
	}

	return id, nil
}
