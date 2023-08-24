package ports

import (
	"errors"
	"gochat/internal/core/domain"
)

var (
	ChannelNotFound error = errors.New("channel not found")
)

type Repositories interface {
	CreateNewChannel(*domain.Channel) error
	CreateIfNotFound(*domain.User) error

	ValidChannel(channelid int) (bool, error)
	UserJoinChannel(userid string, channelid int) error
	ChannelCreatedByUser(userid string, channelid int) (ok bool, err error)
	DeleteChannel(channelid int) error
	PostMessageInChannel(channelid int, m *domain.Message) error

	DeleteIfExistsUser(userid string) error
	GetUserChannels(userid string) ([]domain.ChannelBanner, error)
	GetChannelMessages(channelid int) (*domain.ChannelMessages, error)
}
