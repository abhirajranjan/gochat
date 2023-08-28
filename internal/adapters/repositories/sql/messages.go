package sql

import (
	"context"
	"gochat/internal/core/domain"

	"github.com/pkg/errors"
)

func (r *sqlRepo) deleteMessages(ctx context.Context, cond *Messages) error {
	tx := r.conn.WithContext(ctx)
	if err := tx.Delete(&Messages{}, &cond).Error; err != nil {
		return errors.Wrap(err, "gorm.Delete")
	}
	return nil
}

func (r *sqlRepo) DeleteMessagesByChannelID(ctx context.Context, channelid int) error {
	cond := &Messages{
		ChannelID: channelid,
	}
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	if err := r.deleteMessages(ctx, cond); err != nil {
		return errors.Wrap(err, "deleteMessages")
	}
	return nil
}

func (r *sqlRepo) DeleteMessagesByUserID(ctx context.Context, userid string) error {
	cond := &Messages{
		UserID: userid,
	}
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	if err := r.deleteMessages(ctx, cond); err != nil {
		return errors.Wrap(err, "deleteMessages")
	}
	return nil
}

func (r *sqlRepo) PostMessageInChannel(ctx context.Context, channelid int, m *domain.Message) error {
	message := &Messages{
		UserID:    m.User.ID,
		ChannelID: channelid,
		Content:   m.Content,
		Type:      m.Type,
	}

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if err := tx.Model(&Messages{}).Create(&message).Error; err != nil {
		return errors.Wrap(err, "tx.Create")
	}

	return nil
}

func (r *sqlRepo) GetChannelMessages(ctx context.Context, channelid int) (*domain.ChannelMessages, error) {
	var (
		chanmessages   domain.ChannelMessages
		messageuser    map[string]domain.UserProfile
		domainmessages []domain.Message
		messages       []Messages
		cond           = Messages{
			ChannelID: channelid,
		}
	)

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if err := tx.Where(&cond).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}

	{
		check := map[string]struct{}{}
		messageuserid := make([]string, 0, len(messages))
		for _, m := range messages {
			if _, ok := check[m.UserID]; !ok {
				messageuserid = append(messageuserid, m.UserID)
				check[m.UserID] = struct{}{}
			}
		}
		user, err := r.getUserFromUserID(ctx, messageuserid...)
		if err != nil {
			return nil, err
		}

		for _, u := range user {
			messageuser[u.ID] = domain.UserProfile{
				ID:         u.ID,
				GivenName:  u.GivenName,
				FamilyName: u.FamilyName,
				Picture:    u.Picture,
				NameTag:    u.NameTag,
			}
		}
	}

	domainmessages = make([]domain.Message, len(messages))
	for idx, m := range messages {
		domainmessages[idx] = domain.Message{
			Id:      m.ID,
			User:    messageuser[m.UserID],
			At:      m.CreatedAt,
			Type:    m.Type,
			Content: m.Content,
		}
	}
	chanmessages = domain.ChannelMessages{
		ChannelId: channelid,
		Messages:  domainmessages,
	}

	return &chanmessages, nil
}

func (r *sqlRepo) getRecentMessagesFromChannelID(ctx context.Context, channelid ...int) ([]Messages, error) {
	var messages []Messages
	tx := r.conn.WithContext(ctx)
	window := tx.Raw("Select *, ROW_NUMBER() OVER(PARTITION BY channel_id ORDER BY created_at DESC) AS _rn from messages")

	if err := tx.Table("(?) as t", window).
		Select("t.*").
		Where("t._rn = 1 AND t.channel_id IN ?", channelid).
		Order("t.created_at DESC").
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
