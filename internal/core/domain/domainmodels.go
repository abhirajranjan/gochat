package domain

import "time"

// login request model

type LoginRequest struct {
	Email       string
	Name        string
	Picture     string
	Given_name  string
	Family_name string
}

// user model

type User struct {
	UserID     int64
	GivenName  string
	FamilyName string
	Email      string
	Picture    string
}

// channel models

type MessageType int

const (
	MessageTypeImage = iota + 1
	MessageTypeSticker
	MessageTypeText
)

type Message struct {
	Id      int64
	UserId  int64
	At      time.Time
	Type    MessageType
	Content []byte
}

type Channel struct {
	ChannelBanner
	Users     []User
	CreatedBy User
	Messages  []Message
}

type ChannelBanner struct {
	Id            int64
	Name          string
	Picture       string
	RecentMessage Message
}

type ChannelMessages struct {
	Id       int64 // channel id
	Messages []Message
}
