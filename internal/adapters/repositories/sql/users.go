package sql

import (
	"context"
	"gochat/internal/core/domain"
)

func (r *sqlRepo) NewUser(ctx context.Context, user *domain.User) error {
	var (
		res  User
		cond User = User{
			ID:      user.ID,
			NameTag: user.NameTag,
		}
		attrs = User{
			GivenName:  user.GivenName,
			FamilyName: user.FamilyName,
			Email:      user.Email,
			Picture:    user.Picture,
		}
	)

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)

	if err := tx.Where(&cond).Attrs(&attrs).FirstOrCreate(&res).Error; err != nil {
		return err
	}

	return nil
}

func (r *sqlRepo) DeleteUser(ctx context.Context, userid string) error {
	var cond = User{
		ID: userid,
	}

	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if err := tx.Delete(&cond).Error; err != nil {
		return err
	}

	return nil
}

func (r *sqlRepo) getUserFromUserID(ctx context.Context, userid ...string) ([]User, error) {
	var users []User
	tx := r.conn.WithContext(ctx)

	if err := tx.Model(User{}).
		Where("id IN ?", userid).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *sqlRepo) VerifyUser(ctx context.Context, userid string) (error, bool) {
	var (
		count int64
		cond  User = User{
			ID: userid,
		}
	)
	ctx, cancel := r.getContextWithTimeout(ctx)
	defer cancel()

	tx := r.conn.WithContext(ctx)
	if err := tx.Model(&User{}).Where(&cond).Count(&count).Error; err != nil {
		return err, false
	}
	if count == 0 {
		return nil, false
	}
	return nil, true
}
