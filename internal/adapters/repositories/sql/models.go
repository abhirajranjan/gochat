package sql

import (
	"gochat/internal/core/domain"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GivenName  string `gorm:"column=given_name; not null"`
	FamilyName string `gorm:"column=family_name; not null"`
	Email      string `gorm:"column=email; not null; unique"`
	Picture    string `gorm:"column=picture"`
}

type UserChannels struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID    int64 `gorm:"not null"`
	ChannelID int64 `gorm:"not null"`
}

type Messages struct {
	gorm.Model
	UserID    int64  `gorm:"not null"`
	ChannelID int64  `gorm:"not null"`
	Content   []byte `gorm:"not null"`
	Type      domain.MessageType
}

type Channel struct {
	gorm.Model
	CreatedBy int64
	Name      string
	Picture   string
}
