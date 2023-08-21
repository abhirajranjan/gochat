package services

import (
	"gochat/internal/core/domain"

	"github.com/pkg/errors"
)

func (s *service) JoinChannel(userid string, channelid int) error {
	if err := isValidChannel(channelid, s.repo.ValidChannel); err != nil {
		return err
	}

	ok, err := s.repo.UserinChannel(userid, channelid)
	if err != nil {
		return errors.Wrap(err, "repo.UserinChannel")
	}
	if ok {
		// user already in channel
		return nil
	}

	if err := s.repo.UserJoinChannel(userid, channelid); err != nil {
		return errors.Wrap(err, "repo.UserJoinChannel")
	}

	return nil
}

func (s *service) DeleteChannel(userid string, channelid int) error {
	if err := isValidChannel(channelid, s.repo.ValidChannel); err != nil {
		return err
	}

	ok, err := s.repo.UserinChannel(userid, channelid)
	if err != nil {
		return errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return ErrUserNotInChannel
	}

	if err := s.repo.DeleteChannel(channelid); err != nil {
		return errors.Wrap(err, "repo.DeleteChannel")
	}

	return nil
}

func (s *service) GetMessagesFromChannel(userid string, channelid int) (*domain.ChannelMessages, error) {
	if err := isValidChannel(channelid, s.repo.ValidChannel); err != nil {
		return nil, err
	}

	ok, err := s.repo.UserinChannel(userid, channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return nil, ErrUserNotInChannel
	}

	channelmsg, err := s.repo.GetChannelMessages(channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.GetChannelMessages")
	}

	return channelmsg, nil
}

func (s *service) PostMessageInChannel(userid string, channelid int, message *domain.Message) (*domain.Message, error) {
	if err := isValidChannel(channelid, s.repo.ValidChannel); err != nil {
		return nil, err
	}

	ok, err := s.repo.UserinChannel(userid, channelid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.UserinChannel")
	}
	if !ok {
		return nil, ErrUserNotInChannel
	}

	if len(message.Content) == 0 {
		return nil, ErrCannotBeEmpty("message")
	}

	message.User = domain.UserProfile{ID: userid}

	if err := s.repo.PostMessageInChannel(channelid, message); err != nil {
		return nil, errors.Wrap(err, "PostMessageInChannel")
	}

	return message, nil
}

func isValidChannel(channelid int, validChannel func(channelid int) (bool, error)) error {
	ok, err := validChannel(channelid)
	if err != nil {
		return errors.Wrap(err, "validChannel")
	}
	if !ok {
		return ErrChannelNotFound
	}

	return nil
}
