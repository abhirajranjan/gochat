package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"gochat/internal/core/domain"
)

func (h *handler) AuthMiddleware() func(next http.Handler) http.Handler {
	return h.authMiddleware
}

// injects user information from jwt to context
func (h *handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			jwtoken    string
			session    *sessions.Session
			sessionjwt domain.SessionJwtModel
		)

		// get the session data
		session, err := h.store.Get(r, "session")
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
		}

		IJWTtoken, ok := session.Values["token"]
		if !ok {
			http.Error(w, "", http.StatusBadRequest)
		}

		if jwtoken, ok = IJWTtoken.(string); !ok {
			http.Error(w, "", http.StatusBadRequest)
		}

		_, err = h.jwtParser.ParseWithClaims(jwtoken,
			&sessionjwt,
			func(j *jwt.Token) (interface{}, error) {
				return h.getEncryptionKey(j.Method)
			})

		if err != nil {
			h.Debugf("jwtParser.ParseWithClaims: %s", err)
			switch err {
			case jwt.ErrTokenExpired:
				http.Error(w, "token expired", http.StatusUnauthorized)
				break
			default:
				http.Error(w, "invalid token", http.StatusUnauthorized)
			}
			return
		}

		err, ok = h.service.VerifyUser(sessionjwt.Sub)
		if err != nil {
			h.logger.Debugf("service.VerifyUser: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !ok {
			h.logger.Debugf("service.VerifyUser: user does not exist")
			http.Error(w, "user does not exists", http.StatusBadRequest)
			return
		}

		session.Values[ID_KEY] = sessionjwt.Sub
		next.ServeHTTP(w, r)
	})
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
