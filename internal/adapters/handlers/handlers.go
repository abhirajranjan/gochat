package handlers

import (
	"fmt"
	"gochat/config"
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"gochat/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type handler struct {
	logger  logger.ILogger
	service ports.Service

	wsUpgrader *websocket.Upgrader
	jwtParser  *jwt.Parser

	config config.JwtConfig
}

// handler implements ports.Handler
var _ ports.Handler = (*handler)(nil)

func NewHandler(config config.JwtConfig, s ports.Service, l logger.ILogger) *handler {
	return &handler{
		logger:    l,
		service:   s,
		jwtParser: jwt.NewParser(),
		wsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		config: config,
	}
}

// returns recent user messages
func (h *handler) GetUserMessages(ctx *gin.Context) {
	userid := ctx.GetInt64("userId")
	if userid == 0 {
		h.logger.Errorf("GetUserMessages: %w", errors.Errorf("error getting userid"))
		setInternalServerError(ctx)
		return
	}

	channelbanner, err := h.service.GetUserMessages(userid)
	if err != nil {
		h.logger.Error(err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channelbanner)
}

// extract messages from channels
func (h *handler) GetMessagesFromChannel(ctx *gin.Context) {
	channelid_string := ctx.Params.ByName("channelid")
	if channelid_string == "" {
		h.logger.Debug("GetMessagesFromChannel: no channelid passed")
		setBadRequest(ctx)
		return
	}

	channelid, err := strconv.Atoi(channelid_string)
	if err != nil {
		h.logger.Debugf("GetMessagesFromChannel: strconv.Atoi: %w", err)
		setBadRequest(ctx)
		return
	}

	channelMessages, err := h.service.GetMessagesFromChannel(int64(channelid))
	if err != nil {
		h.logger.Debugf("GetMessagesFromChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channelMessages)
}

// post a new message to channel
func (h *handler) PostMessageInChannel(ctx *gin.Context) {
	// extract user info set by middlewares
	userid := ctx.GetInt64("userId")
	if userid == 0 {
		h.logger.Errorf("GetUserMessages: %w", errors.Errorf("error getting userid"))
		setInternalServerError(ctx)
		return
	}

	channelid_string := ctx.Params.ByName("channelid")
	if channelid_string == "" {
		h.logger.Debug("PostMessageInChannel: no channelid passed")
		setBadRequest(ctx)
		return
	}

	channelid, err := strconv.Atoi(channelid_string)
	if err != nil {
		h.logger.Debugf("PostMessageInChannel: strconv.Atoi: %w", err)
		setBadRequest(ctx)
		return
	}

	var message domain.Message
	if err := ctx.Bind(&message); err != nil {
		h.logger.Debugf("PostMessageInChannel: ctx.Bind: %w", err)
		setBadRequest(ctx)
		return
	}

	if message.UserId != userid {
		h.logger.Debug("PostMessageInChannel: request sender does not match with message embeded user")
		setBadRequestWithErr(ctx, fmt.Errorf("request sender does not match with message embeded user"))
		return
	}

	msg, err := h.service.PostMessageInChannel(int64(channelid), &message)
	if err != nil {
		h.logger.Debugf("PostMessageInChannel: %w", err)
		setInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, msg.Id)
}

func (h *handler) HandleWS(ctx *gin.Context) {
	// upgrade connection to websocket
	conn, err := h.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// in case of error, websocket upgrade automatically
		// write the error to writer with status code
		h.logger.Errorf("ws upgrade: %w", err)
		setInternalServerError(ctx)
	}

	if err := h.service.HandleWS(conn); err != nil {
		h.logger.Errorf("handleWS: %w", err)
		setInternalServerError(ctx)
		return
	}
}
