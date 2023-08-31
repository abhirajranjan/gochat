package handlers

import (
	"gochat/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *handler) HandleGoogleAuth() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cred         Credential
			claims       domain.GoogleJWTModel
			loginRequest domain.LoginRequest
		)

		if err := bindReader(&cred, r.Body, 100000); err != nil {
			h.Debugf("HandleGoogleAuth: bindReader, %s", err)
			http.Error(w, "no token", http.StatusBadRequest)
			return
		}

		h.Debugf("HandleGoogleAuth: %s", cred.Token)

		_, _, err := h.parseUnverified(cred.GetToken(), &claims)
		if err != nil {
			h.Debugf("jwtParser: %w", err)
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		h.Debugf("claims: %#v", claims)

		loginRequest = domain.LoginRequest{
			Email:       claims.Email,
			Name:        claims.Name,
			Picture:     claims.Picture,
			Given_name:  claims.Given_name,
			Family_name: claims.Family_name,
			Sub:         claims.Sub,
		}

		user, err := h.service.NewUser(loginRequest)
		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			h.Debugf("service.LoginRequest: %s", err)
			h.logger.Debugf("")
			setBadRequest(ctx)
			return
		} else if err != nil {
			h.Errorf("service.LoginRequest: %s", err)
			setInternalServerError(ctx)
			return
		}

		jwt, err := h.generateSessionJwt(user)
		if err != nil {
			h.Errorf("token.SignedString: %s", err)
			setInternalServerError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, jwt)
		h.Debugf("response jwt: %s", jwt)
	})
}

// returns recent user messages
func (h *handler) GetUserMessages(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)
	h.logger.Debugf("[%s] called GetUserMessages", userid)

	channelbanner, err := h.service.GetUserMessages(userid)
	if err != nil {
		h.Errorf("GetUserMessages: service.GetUserMessages: ", err)
		setInternalServerError(ctx)
		return
	}

	h.logger.Debugf("res: %#v", channelbanner)
	ctx.JSON(http.StatusOK, channelbanner)
}

func (h *handler) DeleteUser(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	if err := h.service.DeleteUser(userid); err != nil {
		h.Errorf("service.DeleteUser: %s", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
