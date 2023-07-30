package services

import (
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"gochat/logger"

	"github.com/gorilla/websocket"
)

type service struct {
	logger logger.ILogger
}

var _ ports.Service = (*service)(nil)

func NewService(l logger.ILogger) *service {
	return &service{
		logger: l,
	}
}

func (s *service) LoginRequest(c domain.LoginRequest) (domain.User, error)                {}
func (s *service) HandleWS(*websocket.Conn) error                                         {}
func (s *service) GetUserMessages(userId int64) ([]domain.ChannelBanner, error)           {}
func (s *service) GetMessagesFromChannel(channelid int64) (domain.ChannelMessages, error) {}
func (s *service) PostMessageInChannel(channelid int64, message domain.Message) (domain.Message, error) {
}
