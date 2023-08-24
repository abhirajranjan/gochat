package sql

import (
	"context"
	"gochat/internal/core/domain"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *sqlRepo) ValidChannel(channelid int) (bool, error) {
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

func (r *sqlRepo) UserJoinChannel(userid string, channelid int) error {
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	ok, err := r.ValidChannel(channelid)
	if err != nil {
		return err
	}
	if !ok {
		return domain.NewErrDomain("channel doesnt exists")
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

func (r *sqlRepo) UserinChannel(userid string, channelid int) (ok bool, err error) {
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

func (r *sqlRepo) DeleteChannel(channelid int) error {
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

func (r *sqlRepo) PostMessageInChannel(channelid int, m *domain.Message) error {
	message := &Messages{
		UserID:    m.User.ID,
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
	var channelbanner []domain.ChannelBanner

	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)

	banner, err := getUserChannelswithMessage(tx, userid)
	if err != nil {
		return nil, errors.Wrap(err, "getUserChannelswithMessage")
	}

	channelbanner = make([]domain.ChannelBanner, len(banner))
	for idx, b := range banner {
		channelbanner[idx] = domain.ChannelBanner{
			Id:      b.Channel_id,
			Name:    b.Channel_name,
			Picture: b.Channel_picture,
			RecentMessage: domain.Message{
				Id:      b.Message_id,
				At:      b.Message_created_at,
				Type:    domain.MessageType(b.Message_type),
				Content: b.Message_content,
				User: domain.UserProfile{
					ID:         b.User_id,
					GivenName:  b.User_given_name,
					FamilyName: b.User_family_name,
					Picture:    b.User_picture,
					NameTag:    b.User_name_tag,
				},
			},
		}
	}
	return channelbanner, nil
}

func (r *sqlRepo) GetChannelMessages(channelid int) (*domain.ChannelMessages, error) {
	return nil, errors.New("no implemented")
	channelmessages := domain.ChannelMessages{}
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

func (r *sqlRepo) CreateNewChannel(channel *domain.Channel) error {
	ctx, cancel := r.getContextWithTimeout(context.Background())
	defer cancel()

	tx := r.conn.WithContext(ctx)
	querychannel := Channel{
		Name:      channel.Name,
		Desc:      channel.Desc,
		Picture:   channel.Picture,
		CreatedBy: channel.CreatedBy,
	}

	if err := tx.Model(&Channel{}).Create(&querychannel).Error; err != nil {
		return err
	}

	channel.Id = querychannel.ID
	return nil
}

func (r *sqlRepo) getContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, r.config.SqlTimeout)
}
