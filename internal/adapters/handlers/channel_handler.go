package handlers

import (
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

		if userid, err = getUserID(h.store, r); err != nil {
			h.logger.Debugf("NewChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if err := bindReader(&newChanReq, r.Body, 100000); err != nil {
			h.logger.Debugf("NewChannel: decodeReader: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		h.logger.Debugf("query: %v", newChanReq)

		channel, err := h.service.NewChannel(userid, newChanReq)
		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			h.logger.Debugf("service.NewChannel: %s", err)
			setClientBadRequest(w, errdomain.Reason())
			return
		} else if err != nil {
			h.logger.Errorf("service.NewChannel: %s", err)
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

		if userid, err = getUserID(h.store, r); err != nil {
			h.logger.Debugf("DeleteUser: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		channelid, err := getChannelId(r)
		if err != nil {
			h.Debugf("DeleteUser: extractChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		err = h.service.DeleteChannel(userid, channelid)
		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			h.Debugf("service.DeleteChannel: %s", err)
			setClientBadRequest(w, errdomain.Reason())
			return
		} else if err != nil {
			h.Errorf("service.DeleteChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		http.StatusText(http.StatusOK)
	})
}

// extract messages from channels
func (h *handler) GetMessagesFromChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, err := getUserID(h.store, r)
		if err != nil {
			h.logger.Debugf("GetMessagesFromChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		channelid, err := getChannelId(r)
		if err != nil {
			h.Debugf("GetMessagesFromChannel: extractChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return

		}

		channelMessages, err := h.service.GetMessagesFromChannel(userid, channelid)
		if err != nil {
			h.Errorf("GetMessagesFromChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, channelMessages)
	})
}

// post a new message to channel
func (h *handler) PostMessageInChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, err := getUserID(h.store, r)
		if err != nil {
			h.logger.Debugf("PostMessageInChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		channelid, err := getChannelId(r)
		if err != nil {
			h.Debugf("PostMessageInChannel: getChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var msgreq domain.MessageRequest
		if err := bindReader(&msgreq, r.Body, 100000); err != nil {
			h.Debugf("PostMessageInChannel: decodeReader: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		h.logger.Debugf("message: %#v", msgreq)
		msg, err := h.service.NewMessageInChannel(userid, channelid, &msgreq)

		var errdomain domain.ErrDomain
		if errors.As(err, &errdomain) {
			h.Debugf("PostMessageInChannel: service.PostMessageInChannel: %s", errdomain)
			http.Error(w, errdomain.Reason(), http.StatusBadRequest)
			return
		} else if err != nil {
			h.Errorf("PostMessageInChannel: service.PostMessageInChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		setResponseJSON(w, http.StatusOK, msg.Id)
	})
}

func (h *handler) JoinChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userid, err := getUserID(h.store, r)
		if err != nil {
			h.logger.Debugf("JoinChannel: getUserID: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		channelid, err := getChannelId(r)
		if err != nil {
			h.Debugf("JoinChannel: getChannelId: %s", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		err = h.service.JoinChannel(userid, channelid)
		var errdomain domain.ErrDomain
		if errors.Is(err, &errdomain) {
			h.logger.Debug("JoinChannel: service.JoinChannel: %s", errdomain)
			http.Error(w, errdomain.Reason(), http.StatusBadRequest)
			return
		} else if err != nil {
			h.Errorf("JoinChannel: service.JoinChannel: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		http.StatusText(http.StatusOK)
	})
}
