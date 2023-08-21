package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
)

// extract messages from channels
func (h *handler) GetMessagesFromChannel(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId: %w", err)
		setBadRequest(ctx)
		return
	}

	channelMessages, err := h.service.GetMessagesFromChannel(userid, channelid)
	if err != nil {
		h.Errorf("GetMessagesFromChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channelMessages)
}

// post a new message to channel
func (h *handler) PostMessageInChannel(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId: %w", err)
		setBadRequest(ctx)
		return
	}

	var message domain.Message
	if err := ctx.BindJSON(&message); err != nil {
		h.Debugf("ctx.Bind: %w", err)
		return
	}

	if len(message.Content) == 0 {
		h.Debugf("invalid credential provided %#v", message)
		setBadRequest(ctx)
		return
	}

	msg, err := h.service.PostMessageInChannel(userid, channelid, &message)
	if err != nil {
		h.Errorf("PostMessageInChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, msg.Id)
}

func (h *handler) JoinChannel(ctx *gin.Context) {
	// extract user info set by middlewares
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId %w", err)
		setBadRequest(ctx)
		return
	}

	err = h.service.JoinChannel(userid, channelid)
	if errors.Is(err, ports.ErrDomain) {
		setBadRequestWithErr(ctx, err)
		return
	} else if err != nil {
		h.Errorf("JoinChannel: service.JoinChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *handler) DeleteChannel(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId: %w", err)
		setBadRequest(ctx)
		return
	}

	err = h.service.DeleteChannel(userid, channelid)
	if errors.Is(err, ports.ErrDomain) {
		h.Debugf("service.DeleteChannel: %w", err)
		setBadRequestWithErr(ctx, err)
		return
	} else if err != nil {
		h.Errorf("service.DeleteChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *handler) NewChannel(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

func extractChannelId(ctx *gin.Context) (int, error) {
	channelid_string := ctx.Params.ByName("channelid")
	if channelid_string == "" {
		return 0, errors.New("no channelid passed")
	}

	channelid, err := strconv.Atoi(channelid_string)
	if err != nil {
		return 0, errors.Wrap(err, "strconv.Atoi")
	}
	return channelid, nil
}
