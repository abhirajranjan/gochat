package services

import (
	"context"
	"fmt"
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"unicode"

	"github.com/pkg/errors"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func (s *service) NewChannel(userid string, chanreq domain.ChannelRequest) (*domain.Channel, error) {
	if chanreq.Name == "" || !isPrintable(chanreq.Name) {
		return nil, domain.NewErrDomain("invalid name")
	}

	if !isPrintable(chanreq.Desc) {
		return nil, domain.NewErrDomain("invalid desc")
	}

	channel := domain.Channel{
		Name:      chanreq.Name,
		Picture:   chanreq.Picture,
		Desc:      chanreq.Desc,
		CreatedBy: userid,
	}

	ctx := context.Background()

	if err := s.repo.NewChannel(ctx, &channel); err != nil {
		return nil, errors.Wrap(err, "repo.CreateNewChannel")
	}

	if err := s.repo.UserJoinChannel(ctx, userid, channel.Id); err != nil {
		return nil, errors.Wrap(err, "repo.UserJoinChannel")
	}

	if err := s.broadcastMessageToChannel(ctx, channel.Id, domain.BroadcastNewChannel); err != nil {
		return nil, errors.Wrap(err, "broadcastMessageToChannel")
	}

	return &channel, nil
}

func (s *service) DeleteChannel(userid string, channelid int) error {
	ctx := context.Background()

	if err := s.isValidChannel(ctx, channelid); err != nil {
		return errors.Wrap(err, "isValidChannel")
	}

	ok, err := s.repo.IsChannelCreatedByUser(ctx, userid, channelid)
	if err != nil {
		return errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return domain.NewErrDomain("permission denied")
	}

	if err := s.repo.DeleteUserChannelByChannelID(ctx, channelid); err != nil {
		return errors.Wrap(err, "repo.DeleteUserChannelByChannelID")
	}

	if err := s.repo.DeleteChannel(ctx, channelid); err != nil {
		return errors.Wrap(err, "repo.DeleteChannel")
	}

	go func(r ports.Repositories, channelid int) {
		if err := r.DeleteMessagesByChannelID(ctx, channelid); err != nil {
			fmt.Println(errors.Wrap(err, "DeleteMessagesByChannelID"))
		}
	}(s.repo, channelid)

	return nil
}

func (s *service) JoinChannel(userid string, channelid int) error {
	ctx := context.Background()

	if err := s.isValidChannel(ctx, channelid); err != nil {
		return errors.Wrap(err, "isValidChannel")
	}

	ok, err := s.repo.UserinChannel(ctx, userid, channelid)
	if err != nil {
		return errors.Wrap(err, "repo.UserinChannel")
	}
	if ok { // user already in channel
		return nil
	}

	if err := s.repo.UserJoinChannel(ctx, userid, channelid); err != nil {
		return errors.Wrap(err, "repo.UserJoinChannel")
	}

	return nil
}

func (s *service) GetMessagesFromChannel(userid string, channelid int) (*domain.ChannelMessages, error) {
	ctx := context.Background()

	// checks for validity of channel
	if err := s.isValidChannel(ctx, channelid); err != nil {
		return nil, errors.Wrap(err, "isValidChannel")
	}

	// checks if user in channel
	ok, err := s.repo.UserinChannel(ctx, userid, channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return nil, domain.NewErrDomain("permission denied")
	}

	// extract messages from channel
	channelmsg, err := s.repo.GetChannelMessages(ctx, channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.GetChannelMessages")
	}

	return channelmsg, nil
}

func (s *service) NewMessageInChannel(userid string, channelid int, msgreq *domain.MessageRequest) (*domain.Message, error) {
	ctx := context.Background()

	if len(msgreq.Content) == 0 {
		return nil, domain.NewErrDomain("message cannot be empty")
	}

	message := domain.Message{
		User: domain.UserProfile{
			ID: userid,
		},
		Type:    msgreq.Type,
		Content: msgreq.Content,
	}

	ok, err := s.repo.UserinChannel(ctx, userid, channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return nil, domain.NewErrDomain("permission denied")
	}

	if err := s.repo.PostMessageInChannel(ctx, channelid, &message); err != nil {
		return nil, errors.Wrap(err, "PostMessageInChannel")
	}

	return &message, nil
}

func (s *service) broadcastMessageToChannel(ctx context.Context,
	channelid int, m domain.MessageBroadcastType) error {

	message := domain.Message{
		Type:    m.Type,
		Content: m.Content,
	}
	if err := s.repo.PostMessageInChannel(ctx, channelid, &message); err != nil {
		return errors.Wrap(err, "repo.PostMessageInChannel")
	}
	return nil
}

func (s *service) isValidChannel(ctx context.Context, channelid int) error {
	ok, err := s.repo.ValidChannel(ctx, channelid)
	if err != nil {
		return errors.Wrap(err, "repo.ValidChannel")
	}
	if !ok {
		return domain.NewErrDomain("channel does not exist")
	}

	return nil
}
