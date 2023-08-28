package sql

import (
	"context"
	"gochat/internal/core/domain"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *sqlRepo) deleteUserChannel(ctx context.Context, cond *UserChannels) error {
	tx := r.conn.WithContext(ctx)
	if err := tx.Delete(&UserChannels{}, cond).Error; err != nil {
		return errors.Wrap(err, "gorm.Delete")
	}
	return nil
}

func (r *sqlRepo) DeleteUserChannelByUserID(ctx context.Context, userid string) error {
	cond := &UserChannels{
		UserID: userid,
	}
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	if err := r.deleteUserChannel(ctx, cond); err != nil {
		return errors.Wrap(err, "deleteUserChannel")
	}
	return nil
}

func (r *sqlRepo) DeleteUserChannelByChannelID(ctx context.Context, channelid int) error {
	cond := &UserChannels{
		ChannelID: channelid,
	}
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	if err := r.deleteUserChannel(ctx, cond); err != nil {
		return errors.Wrap(err, "deleteUserChannel")
	}
	return nil
}

func (r *sqlRepo) UserJoinChannel(ctx context.Context, userid string, channelid int) error {
	var cond = UserChannels{
		UserID:    userid,
		ChannelID: channelid,
	}

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if err := tx.Create(&cond).Error; err != nil {
		return errors.Wrap(err, "gorm.Create")
	}

	return nil
}

func (r *sqlRepo) UserinChannel(ctx context.Context, userid string, channelid int) (bool, error) {
	var (
		mapUserChannel UserChannels
		cond           = UserChannels{
			UserID:    userid,
			ChannelID: channelid,
		}
	)

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	err := tx.First(&mapUserChannel, &cond).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "tx.First")
	}

	return true, nil
}

func (r *sqlRepo) getChannelIDfromUserID(ctx context.Context, userid ...string) ([]int, error) {
	var userChannelID []int
	tx := r.conn.WithContext(ctx)

	if err := tx.Model(&UserChannels{}).
		Select("channel_id").
		Where("user_id IN ?", userid).
		Find(&userChannelID).Error; err != nil {
		return nil, err
	}

	return userChannelID, nil
}

func (r *sqlRepo) GetUserChannels(ctx context.Context, userid string) ([]domain.ChannelBanner, error) {
	var (
		channelid     []int
		channels      []Channel
		messages      []Messages
		messageuserid []string
		messageusers  []User
		channelbanner []domain.ChannelBanner

		mapperChannel = map[int]Channel{}
		mapperUser    = map[string]User{}
		err           error
	)

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	channelid, err = r.getChannelIDfromUserID(ctx, userid)
	if len(channelid) == 0 {
		return []domain.ChannelBanner{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "getChannelIDfromUserIDs")
	}

	channels, err = r.getChannelfromChannelID(ctx, channelid...)
	if err != nil {
		return nil, errors.Wrap(err, "getChannelfromChannelIDs")
	}

	for _, c := range channels {
		mapperChannel[c.ID] = c
	}

	messages, err = r.getRecentMessagesFromChannelID(ctx, channelid...)
	if err != nil {
		return nil, errors.Wrap(err, "getRecentMessagesFromChannelIDs")
	}

	{
		check := make(map[string]struct{})
		messageuserid = make([]string, len(messages))
		for idx, m := range messages {
			if _, ok := check[m.UserID]; !ok {
				messageuserid[idx] = m.UserID
				check[m.UserID] = struct{}{}
			}
		}
	}

	messageusers, err = r.getUserFromUserID(ctx, messageuserid...)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "getUserFromUserIDs")
	}
	for _, u := range messageusers {
		mapperUser[u.ID] = u
	}

	channelbanner = make([]domain.ChannelBanner, len(channels))
	for idx, m := range messages {
		var user domain.UserProfile

		u, ok := mapperUser[m.UserID]
		if ok {
			user = domain.UserProfile{
				ID:         u.ID,
				GivenName:  u.GivenName,
				FamilyName: u.FamilyName,
				Picture:    u.Picture,
				NameTag:    u.NameTag,
			}
		}

		recentmessage := domain.Message{
			Id:      m.ID,
			Type:    m.Type,
			At:      m.CreatedAt,
			Content: m.Content,
			User:    user,
		}

		channel := mapperChannel[m.ChannelID]
		channelbanner[idx] = domain.ChannelBanner{
			RecentMessage: recentmessage,
			Id:            channel.ID,
			Name:          channel.Name,
			Picture:       channel.Picture,
		}
	}

	return channelbanner, nil
}
