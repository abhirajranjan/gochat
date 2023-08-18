package sql

import (
	"gochat/internal/core/domain"
	"time"
)

type User struct {
	ID         string `gorm:"primary key"`
	NameTag    string `gorm:"primary key"`
	GivenName  string `gorm:"column=given_name; not null"`
	FamilyName string `gorm:"column=family_name; not null"`
	Email      string `gorm:"column=email; not null; unique"`
	Picture    string `gorm:"column=picture"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserChannels struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID    string `gorm:"not null"`
	ChannelID int64  `gorm:"not null"`
}

type Messages struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID    string `gorm:"not null"`
	ChannelID int64  `gorm:"not null"`
	Content   []byte `gorm:"not null"`
	Type      domain.MessageType
}

type Channel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int64
	Name      string
	Picture   string
}
