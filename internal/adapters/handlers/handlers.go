package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"gochat/config"
	"gochat/internal/core/ports"
)

const ID_KEY = "NAMETAG"
const JWT_ISSUER = "connector/auth"

type handler struct {
	service ports.Service

	wsUpgrader *websocket.Upgrader
	jwtParser  *jwt.Parser
	config     config.JwtConfig
}

// handler implements ports.Handler
var _ ports.Handler = (*handler)(nil)

func NewHandler(config config.JwtConfig, s ports.Service) *handler {
	return &handler{
		service:   s,
		jwtParser: jwt.NewParser(),
		wsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		config: config,
	}
}
