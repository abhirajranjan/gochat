package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"gochat/config"
	"gochat/internal/core/ports"
)

const NAMETAGKEY = "NAMETAG"
const JWT_ISSUER = "connector/auth"

type logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type handler struct {
	logger  logger
	service ports.Service

	wsUpgrader *websocket.Upgrader
	jwtParser  *jwt.Parser

	config config.JwtConfig
}

// handler implements ports.Handler
var _ ports.Handler = (*handler)(nil)

func NewHandler(config config.JwtConfig, s ports.Service, l logger) *handler {
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

func (h *handler) Debug(args ...interface{}) {
	if h.logger != nil {
		h.logger.Debug(args...)
	}
}

func (h *handler) Debugf(template string, args ...interface{}) {
	if h.logger != nil {
		h.logger.Debugf(template, args...)
	}
}

func (h handler) Errorf(template string, args ...interface{}) {
	if h.logger != nil {
		h.logger.Errorf(template, args...)
	}
}
