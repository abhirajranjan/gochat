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

	JoinChannel(userid string, channelid int64) error
	DeleteChannel(userid string, channelid int64) error
	GetMessagesFromChannel(userid string, channelid int64) (*domain.ChannelMessages, error)
	PostMessageInChannel(userid string, channelid int64, message *domain.Message) (*domain.Message, error)
}
