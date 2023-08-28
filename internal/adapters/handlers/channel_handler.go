package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"gochat/internal/core/domain"
)

func (h *handler) NewChannel(ctx *gin.Context) {
	var newchanreq domain.ChannelRequest
	userid := ctx.GetString(NAMETAGKEY)

	if err := ctx.Bind(&newchanreq); err != nil {
		h.logger.Debugf("ctx.Bind: %w", err)
		setBadRequest(ctx)
		return
	}

	h.logger.Debugf("query: %#v", newchanreq)

	channel, err := h.service.NewChannel(userid, newchanreq)
	var errdomain domain.ErrDomain
	if errors.As(err, &errdomain) {
		h.logger.Debugf("service.NewChannel: %s", err)
		setBadReqWithClientErr(ctx, errdomain)
		return
	} else if err != nil {
		h.logger.Errorf("service.NewChannel: %s", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (h *handler) DeleteChannel(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId: %s", err)
		setBadRequest(ctx)
		return
	}

	err = h.service.DeleteChannel(userid, channelid)
	var errdomain domain.ErrDomain
	if errors.As(err, &errdomain) {
		h.Debugf("service.DeleteChannel: %s", err)
		setBadReqWithClientErr(ctx, errdomain)
		return
	} else if err != nil {
		h.Errorf("service.DeleteChannel: %s", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// extract messages from channels
func (h *handler) GetMessagesFromChannel(ctx *gin.Context) {
	userid := ctx.GetString(NAMETAGKEY)

	channelid, err := extractChannelId(ctx)
	if err != nil {
		h.Debugf("extractChannelId: %s", err)
		setBadRequest(ctx)
		return
	}

	channelMessages, err := h.service.GetMessagesFromChannel(userid, channelid)
	if err != nil {
		h.Errorf("GetMessagesFromChannel: %s", err)
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
		h.Debugf("extractChannelId: %s", err)
		setBadRequest(ctx)
		return
	}

	var msgreq domain.MessageRequest
	if err := ctx.Bind(&msgreq); err != nil {
		h.Debugf("ctx.Bind: %s", err)
		return
	}

	h.logger.Debugf("message: %#v", msgreq)

	msg, err := h.service.NewMessageInChannel(userid, channelid, &msgreq)

	var errdomain domain.ErrDomain
	if errors.As(err, &errdomain) {
		h.Debugf("service.PostMessageInChannel: %s", errdomain)
		return
	} else if err != nil {
		h.Errorf("service.PostMessageInChannel: %s", err)
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
		h.Debugf("extractChannelId %s", err)
		setBadRequest(ctx)
		return
	}

	err = h.service.JoinChannel(userid, channelid)
	var errdomain domain.ErrDomain
	if errors.Is(err, &errdomain) {
		setBadReqWithClientErr(ctx, errdomain)
		return
	} else if err != nil {
		h.Errorf("JoinChannel: service.JoinChannel: %s", err)
		setInternalServerError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
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
