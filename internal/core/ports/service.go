package ports

import (
	"gochat/internal/core/domain"

	"github.com/gorilla/websocket"
)

const USERID_SEP = "@"

type Service interface {
	HandleWS(*websocket.Conn) error

	NewUser(domain.LoginRequest) (*domain.User, error)
	DeleteUser(userid string) error
	VerifyUser(userid string) (error, bool)
	GetUserMessages(userid string) ([]domain.ChannelBanner, error)

	NewChannel(userid string, chanreq domain.ChannelRequest) (*domain.Channel, error)
	DeleteChannel(userid string, channelid int) error
	JoinChannel(userid string, channelid int) error
	NewMessageInChannel(userid string, channelid int, message *domain.MessageRequest) (*domain.Message, error)
	GetMessagesFromChannel(userid string, channelid int) (*domain.ChannelMessages, error)
}
