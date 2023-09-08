package handlers

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"gochat/internal/core/domain"
)

func (h *handler) NewChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			newChanReq domain.ChannelRequest
			userid     string
			err        error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("NewChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if err := bindReader(&newChanReq, r.Body, 100000); err != nil {
			slog.Error("NewChannel: decodeReader: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("query: %v", newChanReq)

		channel, err := h.service.NewChannel(userid, newChanReq)
		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			slog.Error("service.NewChannel: %s", err)
			setClientBadRequest(w, errdomain.Reason())
			return
		}

		if err != nil {
			slog.Error("service.NewChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, channel)
	})
}

func (h *handler) DeleteChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid string
			err    error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("DeleteUser: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		channelid, err := getChannelId(r)
		if err != nil {
			slog.Error("DeleteUser: extractChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("delete Channelid", channelid)

		err = h.service.DeleteChannel(userid, channelid)
		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			slog.Error("service.DeleteChannel: %s", err)
			setClientBadRequest(w, errdomain.Reason())
			return
		}

		if err != nil {
			slog.Error("service.DeleteChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// extract messages from channels
func (h *handler) GetMessagesFromChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid    string
			channelid int
			err       error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Debug("GetMessagesFromChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if channelid, err = getChannelId(r); err != nil {
			slog.Debug("GetMessagesFromChannel: extractChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return

		}

		slog.Debug("GetMessagesFromChannel", "channelid", channelid, "userid", userid)

		channelMessages, err := h.service.GetMessagesFromChannel(userid, channelid)
		if err != nil {
			slog.Error("GetMessagesFromChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, channelMessages)
	})
}

// post a new message to channel
func (h *handler) PostMessageInChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid    string
			channelid int
			err       error
		)

		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("PostMessageInChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if channelid, err = getChannelId(r); err != nil {
			slog.Error("PostMessageInChannel: getChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var msgreq domain.MessageRequest
		if err := bindReader(&msgreq, r.Body, 100000); err != nil {
			slog.Error("PostMessageInChannel: decodeReader: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("message: %#v", msgreq)
		msg, err := h.service.NewMessageInChannel(userid, channelid, &msgreq)

		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			slog.Error("PostMessageInChannel: service.PostMessageInChannel: %s", errdomain)
			http.Error(w, errdomain.Reason(), http.StatusBadRequest)
			return
		}

		if err != nil {
			slog.Error("PostMessageInChannel: service.PostMessageInChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, msg.Id)
	})
}

func (h *handler) JoinChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userid    string
			channelid int
			err       error
		)
		if userid, err = getUserID(r.Context()); err != nil {
			slog.Error("JoinChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if channelid, err = getChannelId(r); err != nil {
			slog.Error("JoinChannel: getChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		err = h.service.JoinChannel(userid, channelid)

		var errdomain domain.ErrDomain
		if errors.Is(err, &errdomain) {
			slog.Error("JoinChannel: service.JoinChannel: %s", errdomain)
			http.Error(w, errdomain.Reason(), http.StatusBadRequest)
			return
		}

		if err != nil {
			slog.Error("JoinChannel: service.JoinChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
