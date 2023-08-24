package ports

import (
	"gochat/internal/core/domain"

	"github.com/gorilla/websocket"
)

const USERID_SEP = "@"

type Service interface {
	HandleWS(*websocket.Conn) error

	LoginRequest(domain.LoginRequest) (*domain.User, error)
	DeleteUser(userid string) error
	GetUserMessages(userid string) ([]domain.ChannelBanner, error)

	JoinChannel(userid string, channelid int) error
	NewChannel(userid string, chanreq domain.NewChannelRequest) (*domain.Channel, error)
	DeleteChannel(userid string, channelid int) error
	GetMessagesFromChannel(userid string, channelid int) (*domain.ChannelMessages, error)
	PostMessageInChannel(userid string, channelid int, message *domain.Message) (*domain.Message, error)
}
