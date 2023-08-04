package ports

import "gochat/internal/core/domain"

type Repositories interface {
	CreateIfNotFound(*domain.User) error
	ValidUser(userid int64) (bool, error)
	GetUserChannels(userid int64) ([]domain.ChannelBanner, error)
	ValidChannel(channelid int64) (bool, error)
	GetChannelMessages(channelid int64) (*domain.ChannelMessages, error)
	PostMessageInChannel(channelid int64, m *domain.Message) error
}
