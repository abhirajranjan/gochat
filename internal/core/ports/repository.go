package ports

import (
	"context"
	"errors"
	"gochat/internal/core/domain"
)

var (
	ChannelNotFound error = errors.New("channel not found")
)

type Repositories interface {
	NewUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userid string) error
	VerifyUser(ctx context.Context, userid string) (error, bool)

	NewChannel(ctx context.Context, channel *domain.Channel) error
	DeleteChannel(ctx context.Context, channelid int) error
	IsChannelCreatedByUser(ctx context.Context, userid string, channelid int) (ok bool, err error)
	ValidChannel(ctx context.Context, channelid int) (bool, error)

	DeleteUserChannelByChannelID(ctx context.Context, channelid int) error
	DeleteUserChannelByUserID(ctx context.Context, userid string) error
	UserinChannel(ctx context.Context, userid string, channelid int) (bool, error)
	UserJoinChannel(ctx context.Context, userid string, channelid int) error
	GetUserChannels(ctx context.Context, userid string) ([]domain.ChannelBanner, error)

	PostMessageInChannel(ctx context.Context, channelid int, m *domain.Message) error
	GetChannelMessages(ctx context.Context, channelid int) (*domain.ChannelMessages, error)
	DeleteMessagesByChannelID(ctx context.Context, channelid int) error
	DeleteMessagesByUserID(ctx context.Context, userid string) error
}
