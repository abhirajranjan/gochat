package services

import (
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"

	"github.com/pkg/errors"
)

func (s *service) LoginRequest(c domain.LoginRequest) (*domain.User, error) {
	var user domain.User

	if c.Email == "" {
		return nil, ErrCannotBeEmpty("email")
	}

	if c.Family_name == "" {
		return nil, ErrCannotBeEmpty("family name")
	}

	if c.Given_name == "" {
		return nil, ErrCannotBeEmpty("given name")
	}
	if c.Name == "" {
		return nil, ErrCannotBeEmpty("name")
	}

	if c.Sub == "" {
		return nil, ErrCannotBeEmpty("Sub")
	}

	user.ID = c.Sub
	user.NameTag = generateNameTag(c.Given_name, c.Sub)
	user.Email = c.Email
	user.GivenName = c.Given_name
	user.FamilyName = c.Family_name
	user.Picture = c.Picture

	if err := s.repo.CreateIfNotFound(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) DeleteUser(userid string) error {
	if userid == "" {
		return ErrUserNotFound
	}

	if err := s.repo.DeleteIfExistsUser(userid); err != nil {
		return errors.Wrap(err, "repo.DeleteIfExistsUser")
	}

	return nil
}

func (s *service) GetUserMessages(userid string) ([]domain.ChannelBanner, error) {
	channelbanner, err := s.repo.GetUserChannels(userid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.GetUserChannels")
	}

	return channelbanner, nil
}

func generateNameTag(name, tag string) string {
	return name + ports.USERID_SEP + tag
}
