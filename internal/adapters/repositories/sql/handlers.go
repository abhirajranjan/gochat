package sql

import (
	"context"
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *sqlRepo) ValidChannel(channelid int64) (bool, error) {
	var (
		count int64
	)

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	if res := r.conn.WithContext(ctx).
		Model(&Channel{}).
		Where(&Channel{
			ID: channelid,
		}).
		Count(&count); res.Error != nil {
		return false, res.Error
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r *sqlRepo) UserJoinChannel(userid string, channelid int64) error {
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	ok, err := r.ValidChannel(channelid)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Wrap(ports.ErrDomain, "channel doesnt exists")
	}

	if res := r.conn.WithContext(ctx).
		Model(&UserChannels{}).
		Create(UserChannels{
			UserID:    userid,
			ChannelID: channelid,
		}); res.Error != nil {
		return errors.Wrap(res.Error, "gorm.Create")
	}

	return nil
}

func (r *sqlRepo) UserinChannel(userid string, channelid int64) (ok bool, err error) {
	var userchan UserChannels

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	res := tx.First(&userchan, UserChannels{UserID: userid, ChannelID: channelid})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if res.Error != nil {
		return false, errors.Wrap(res.Error, "db.First")
	}

	if userchan.UserID != "" {
		return true, nil
	}

	return false, errors.New("record found yet empty")
}

func (r *sqlRepo) DeleteChannel(channelid int64) error {
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if res := tx.
		Model(&UserChannels{}).
		Delete(&UserChannels{
			ChannelID: channelid,
		}); res.Error != nil {
		return errors.Wrap(res.Error, "gorm.Delete")
	}

	if res := tx.Model(&Channel{}).
		Delete(&Channel{
			ID: channelid,
		}); res.Error != nil {
		return errors.Wrap(res.Error, "gorm.Delete")
	}

	return nil
}

func (r *sqlRepo) PostMessageInChannel(channelid int64, m *domain.Message) error {
	message := &Messages{
		UserID:    m.UserId,
		ChannelID: channelid,
		Content:   m.Content,
		Type:      m.Type,
	}

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if res := tx.Model(&Messages{}).Create(&message); res.Error != nil {
		return errors.Wrap(res.Error, "gorm.Create")
	}

	return nil
}

func (r *sqlRepo) CreateIfNotFound(user *domain.User) error {
	var u User

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tx := r.conn.WithContext(ctx)

	res := tx.Where(&User{
		ID:      user.ID,
		NameTag: user.NameTag,
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

func (r *sqlRepo) DeleteIfExistsUser(userid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	res := tx.Delete(&User{
		ID: userid,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *sqlRepo) GetUserChannels(userid string) ([]domain.ChannelBanner, error) {
	var arrChannelBanner []domain.ChannelBanner
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	res := tx.Model(&UserChannels{}).
		Joins("JOINS channel AS c ON c.id = user_channels.channel_id AND user.id = ?", userid).
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

func (r *sqlRepo) getContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, r.config.SqlTimeout)
}
