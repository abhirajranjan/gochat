package handlers

import (
	"gochat/internal/core/domain"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
)

func (h *handler) HandleGoogleAuth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cred         Credential
			claims       domain.GoogleJWTModel
			loginRequest domain.LoginRequest
		)

		if err := bindReader(&cred, r.Body, 100000); err != nil {
			slog.Error("HandleGoogleAuth: bindReader, %s", err)
			http.Error(w, "no token", http.StatusBadRequest)
			return
		}

		slog.Debug("HandleGoogleAuth: %s", cred.Token)

		if _, _, err := h.parseUnverified(cred.GetToken(), &claims); err != nil {
			slog.Error("jwtParser: %w", err)
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		slog.Debug("claims: %#v", claims)

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
			slog.Error("service.LoginRequest: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if err != nil {
			slog.Error("service.LoginRequest: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		jwt, err := h.generateSessionJwt(user)
		if err != nil {
			slog.Error("token.SignedString: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, jwt)
	})
}

// returns recent user messages
func (h *handler) GetUserMessages() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid string
			err    error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("GetUserMessages: userid not found")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		channelbanner, err := h.service.GetUserMessages(userid)
		if err != nil {
			slog.Error("GetUserMessages: service.GetUserMessages: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, channelbanner)
	})
}

func (h *handler) DeleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid string
			err    error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("DeleteUser: getUserID: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := h.service.DeleteUser(userid); err != nil {
			slog.Error("service.DeleteUser: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
