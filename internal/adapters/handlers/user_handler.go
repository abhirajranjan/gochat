package handlers

import (
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *handler) HandleGoogleAuth(ctx *gin.Context) {
	var (
		cred         Credential
		claims       domain.GoogleJWTModel
		loginRequest domain.LoginRequest
	)

	if err := ctx.Bind(&cred); err != nil {
		h.Debugf("ctx.Bind, %w", err)
		return
	}

	if cred.Token == "" {
		setInvalidToken(ctx)
		return
	}

	h.Debugf("HandleGoogleAuth %s", cred.Token)

	_, _, err := h.parseUnverified(cred.GetToken(), &claims)
	if err != nil {
		h.Debugf("jwtParser: %w", err)
		setInvalidToken(ctx)
		return
	}

	loginRequest = domain.LoginRequest{
		Email:       claims.Email,
		Name:        claims.Name,
		Picture:     claims.Picture,
		Given_name:  claims.Given_name,
		Family_name: claims.Family_name,
		Sub:         claims.Sub,
	}

	user, err := h.service.LoginRequest(loginRequest)
	if errors.Is(err, ports.ErrDomain) {
		setBadRequestWithErr(ctx, err)
		return
	} else if err != nil {
		h.Errorf("service.LoginRequest: %w", err)
		setInternalServerError(ctx)
		return
	}

	jwt, err := h.generateSessionJwt(user)
	if err != nil {
		h.Errorf("token.SignedString: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
	h.Debugf("response jwt: %s", jwt)
}

// returns recent user messages
func (h *handler) GetUserMessages(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelbanner, err := h.service.GetUserMessages(userid)
	if err != nil {
		h.Errorf("GetUserMessages: service.GetUserMessages: ", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channelbanner)
}

func (h *handler) DeleteUser(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	if err := h.service.DeleteUser(userid); err != nil {
		h.Errorf("service.DeleteUser: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
