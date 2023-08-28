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
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserChannels struct {
	UserID    string `gorm:"not null"`
	ChannelID int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Messages struct {
	ID        int    `gorm:"primarykey,auto increment"`
	UserID    string `gorm:"not null"`
	ChannelID int    `gorm:"not null"`
	Content   []byte `gorm:"not null"`
	Type      domain.MessageType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Channel struct {
	ID        int `gorm:"primarykey,auto increment"`
	Name      string
	Picture   string
	CreatedBy string
	Desc      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
