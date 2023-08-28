package sql

import (
	"context"
	"gochat/internal/core/domain"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *sqlRepo) NewChannel(ctx context.Context, channel *domain.Channel) error {
	querychannel := Channel{
		Name:      channel.Name,
		Desc:      channel.Desc,
		Picture:   channel.Picture,
		CreatedBy: channel.CreatedBy,
	}

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)

	if err := tx.Model(&Channel{}).Create(&querychannel).Error; err != nil {
		return err
	}

	channel.Id = querychannel.ID
	return nil
}

func (r *sqlRepo) DeleteChannel(ctx context.Context, channelid int) error {
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)

	cond := &Channel{
		ID: channelid,
	}

	if err := tx.Delete(&Channel{}, cond).Error; err != nil {
		return errors.Wrap(err, "gorm.Delete")
	}

	return nil
}

func (r *sqlRepo) IsChannelCreatedByUser(ctx context.Context, userid string, channelid int) (ok bool, err error) {
	var (
		channel Channel
		cond    = Channel{
			CreatedBy: userid,
			ID:        channelid,
		}
	)

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	err = tx.First(&channel, cond).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "db.First")
	}

	return true, nil
}

func (r *sqlRepo) ValidChannel(ctx context.Context, channelid int) (bool, error) {
	var (
		count int64
		cond  = Channel{
			ID: channelid,
		}
	)

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	err := r.conn.WithContext(ctx).
		Model(&Channel{}).
		Where(&cond).
		Count(&count).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (r *sqlRepo) getChannelfromChannelID(ctx context.Context, channelid ...int) ([]Channel, error) {
	var userChannels []Channel
	tx := r.conn.WithContext(ctx)

	if err := tx.Model(&Channel{}).
		Where("id IN ?", channelid).
		Find(&userChannels).Error; err != nil {
		return nil, err
	}

	return userChannels, nil
}
