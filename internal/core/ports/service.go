package ports

import (
	"gochat/internal/core/domain"

	"github.com/gorilla/websocket"
)

type Service interface {
	LoginRequest(c domain.LoginRequest) (domain.User, error)
	HandleWS(*websocket.Conn) error
	GetUserMessages(userId int64) ([]domain.ChannelBanner, error)
	GetMessagesFromChannel(channelid int64) (domain.ChannelMessages, error)
	PostMessageInChannel(channelid int64, message domain.Message) (domain.Message, error)
}
