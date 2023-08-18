package ports

import "gochat/internal/core/domain"

type Repositories interface {
	ValidChannel(channelid int64) (bool, error)
	UserJoinChannel(userid string, channelid int64) error
	UserinChannel(userid string, channelid int64) (ok bool, err error)
	DeleteChannel(channelid int64) error
	PostMessageInChannel(channelid int64, m *domain.Message) error

	CreateIfNotFound(*domain.User) error
	DeleteIfExistsUser(userid string) error
	GetUserChannels(userid string) ([]domain.ChannelBanner, error)
	GetChannelMessages(channelid int64) (*domain.ChannelMessages, error)
}
