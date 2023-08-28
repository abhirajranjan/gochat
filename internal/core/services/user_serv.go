package services

import (
	"context"
	"fmt"
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"

	"github.com/pkg/errors"
)

func (s *service) NewUser(c domain.LoginRequest) (*domain.User, error) {
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

	ctx := context.Background()
	if err := s.repo.NewUser(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) DeleteUser(userid string) error {
	if userid == "" {
		return ErrUserNotFound
	}

	ctx := context.Background()
	if err := s.repo.DeleteUser(ctx, userid); err != nil {
		return errors.Wrap(err, "repo.DeleteUser")
	}

	go func(r ports.Repositories, userid string) {
		ctx := context.Background()
		if err := s.repo.DeleteUserChannelByUserID(ctx, userid); err != nil {
			fmt.Println(errors.Wrap(err, "repo.DeleteUserChannelByUserID"))
		}
	}(s.repo, userid)

	go func(r ports.Repositories, userid string) {
		ctx := context.Background()
		if err := r.DeleteMessagesByUserID(ctx, userid); err != nil {
			fmt.Println(errors.Wrap(err, "repo.DeleteMessagesByUserID"))
		}
	}(s.repo, userid)

	return nil
}

func (s *service) GetUserMessages(userid string) ([]domain.ChannelBanner, error) {
	ctx := context.Background()

	channelbanner, err := s.repo.GetUserChannels(ctx, userid)
	if err != nil {
		return nil, errors.Wrap(err, "repo.GetUserChannels")
	}

	return channelbanner, nil
}

func generateNameTag(name, tag string) string {
	return name + ports.USERID_SEP + tag
}

func (s *service) VerifyUser(userid string) (error, bool) {
	if userid == "" {
		return nil, false
	}

	ctx := context.Background()
	err, ok := s.repo.VerifyUser(ctx, userid)
	if err != nil {
		return err, false
	}
	return nil, ok
}
