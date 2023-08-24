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

// new channel request model
type NewChannelRequest struct {
	Name    string `json:"name" form:"name"`
	Picture string `json:"picture" form:"picture"`
	Desc    string `json:"desc" form:"desc"`
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
	Id      int         `json:"id"`
	User    UserProfile `json:"user"`
	At      time.Time   `json:"at"`
	Type    MessageType `json:"type"`
	Content []byte      `json:",string"`
}

type Channel struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Picture   string `json:"picture"`
	CreatedBy string `json:"created_by"`
}

type ChannelBanner struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Picture       string  `json:"picture"`
	RecentMessage Message `json:"message"`
}

type ChannelMessages struct {
	ChannelId int    `json:"channel_id"`
	Nextptr   string `json:"next_ptr"`
	Messages  []Message
}
