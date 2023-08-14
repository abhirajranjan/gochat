package sql

import (
	"context"
	"gochat/internal/core/domain"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *sqlRepo) CreateIfNotFound(user *domain.User) error {
	var u User

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tx := r.conn.WithContext(ctx)

	res := tx.Where(&User{
		Model: gorm.Model{
			ID: uint(user.UserID),
		},
	}).Attrs(&User{
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Email:      user.Email,
		Picture:    user.Picture,
	}).FirstOrCreate(&u)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *sqlRepo) ValidUser(userid int64) (bool, error) {
	var user User

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if res := tx.First(&user, userid); res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func (r *sqlRepo) GetUserChannels(userid int64) ([]domain.ChannelBanner, error) {
	var arrChannelBanner []domain.ChannelBanner
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	res := tx.Model(&UserChannels{}).
		Joins("JOINS channel AS c ON c.id = user_channels.channel_id AND user.user_id = ?", userid).
		Select("c.id, c.name, c.picture")

	rows, err := res.Rows()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	arrChannelBanner = make([]domain.ChannelBanner, res.RowsAffected)
	{
		i := 0
		for rows.Next() {
			if err := rows.Scan(&arrChannelBanner[i]); err != nil {
				return nil, errors.Wrap(err, "GetUserChannels: rows.Scan")
			}
			i++
		}
	}

	return arrChannelBanner, nil
}

func (r *sqlRepo) ValidChannel(channelid int64) (bool, error) {
	var channel Channel

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if res := tx.First(&channel, channelid); res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func (r *sqlRepo) GetChannelMessages(channelid int64) (*domain.ChannelMessages, error) {
	channelmessages := domain.ChannelMessages{
		Id: channelid,
	}
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	res := r.conn.WithContext(ctx).Where(&Messages{
		ChannelID: channelid,
	}).Order("created_at DESC")

	if res.Error != nil {
		return nil, res.Error
	}

	rows, err := res.Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}
		channelmessages.Messages = append(channelmessages.Messages, m)
	}

	return &channelmessages, nil
}

func (r *sqlRepo) PostMessageInChannel(channelid int64, m *domain.Message) error {
	return nil
}

func (r *sqlRepo) getContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, r.config.SqlTimeout)
}
