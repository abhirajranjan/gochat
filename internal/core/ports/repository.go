package ports

import "gochat/internal/core/domain"

type Repositories interface {
	CreateNewChannel(*domain.Channel) error
	CreateIfNotFound(*domain.User) error

	ValidChannel(channelid int) (bool, error)
	UserJoinChannel(userid string, channelid int) error
	UserinChannel(userid string, channelid int) (ok bool, err error)
	DeleteChannel(channelid int) error
	PostMessageInChannel(channelid int, m *domain.Message) error

	DeleteIfExistsUser(userid string) error
	GetUserChannels(userid string) ([]domain.ChannelBanner, error)
	GetChannelMessages(channelid int) (*domain.ChannelMessages, error)
}
