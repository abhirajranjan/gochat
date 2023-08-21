package domain

import "time"

// login request model

type LoginRequest struct {
	Email       string
	Name        string
	Picture     string
	Given_name  string
	Family_name string
	Sub         string
}

// user model

type UserProfile struct {
	ID         string
	GivenName  string
	FamilyName string
	Picture    string
	NameTag    string
}

type User struct {
	ID         string
	GivenName  string
	FamilyName string
	Email      string
	Picture    string
	NameTag    string
}

// channel models

type MessageType int

const (
	MessageTypeImage = iota + 1
	MessageTypeSticker
	MessageTypeText
)

type Message struct {
	Id      int
	User    UserProfile
	At      time.Time
	Type    MessageType
	Content []byte `json:",string"`
}

type Channel struct {
	Id        int64
	Name      string
	Picture   string
	Users     []User
	CreatedBy User
	Messages  []Message
}

type ChannelBanner struct {
	Id            int
	Name          string
	Picture       string
	RecentMessage Message
}

type ChannelMessages struct {
	Id       int // channel id
	Messages []Message
}
