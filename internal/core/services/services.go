package services

import (
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type service struct {
	repo ports.Repositories
}

var _ ports.Service = (*service)(nil)

func NewService(repo ports.Repositories) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) LoginRequest(c domain.LoginRequest) (*domain.User, error) {
	var user domain.User

	if c.Email == "" {
		return nil, errors.Wrap(ports.ErrDomain, "email cannot be empty")
	}

	if c.Family_name == "" {
		return nil, errors.Wrap(ports.ErrDomain, "family name cannot be empty")
	}

	if c.Given_name == "" {
		return nil, errors.Wrap(ports.ErrDomain, "given name cannot be empty")
	}
	if c.Name == "" {
		return nil, errors.Wrap(ports.ErrDomain, "name cannot be empty")
	}

	user.Email = c.Email
	user.GivenName = c.Given_name
	user.FamilyName = c.Family_name
	user.Picture = c.Picture

	if err := s.repo.CreateIfNotFound(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) HandleWS(*websocket.Conn) error {
	// TODO: add websocket service
	return nil
}

func (s *service) GetUserMessages(userid int64) ([]domain.ChannelBanner, error) {
	ok, err := s.repo.ValidUser(userid)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserMessages")
	}
	if !ok {
		return nil, errors.Wrap(ports.ErrDomain, "invalid userid")
	}

	channelbanner, err := s.repo.GetUserChannels(userid)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserMessages")
	}

	return channelbanner, nil
}

func (s *service) GetMessagesFromChannel(channelid int64) (*domain.ChannelMessages, error) {
	ok, err := s.repo.ValidChannel(channelid)
	if err != nil {
		return nil, errors.Wrap(err, "GetMessagesFromChannel")
	}
	if !ok {
		return nil, errors.Wrap(ports.ErrDomain, "invalid channel")
	}

	channelmsg, err := s.repo.GetChannelMessages(channelid)
	if err != nil {
		return nil, errors.Wrap(err, "GetMessagesFromChannel")
	}

	return channelmsg, nil
}

func (s *service) PostMessageInChannel(channelid int64, message *domain.Message) (*domain.Message, error) {
	ok, err := s.repo.ValidUser(message.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "PostMessageInChannel")
	}
	if !ok {
		return nil, errors.Wrap(ports.ErrDomain, "invalid user")
	}

	if len(message.Content) == 0 {
		return nil, errors.Wrap(ports.ErrDomain, "message cannot be empty")
	}

	if err := s.repo.PostMessageInChannel(channelid, message); err != nil {
		return nil, errors.Wrap(err, "PostMessageInChannel")
	}

	return message, nil
}
