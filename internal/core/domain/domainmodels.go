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
type ChannelRequest struct {
	Name    string `json:"name" form:"name"`
	Picture string `json:"picture" form:"picture"`
	Desc    string `json:"desc" form:"desc"`
}

type MessageRequest struct {
	Type    MessageType `json:"type" form:"type"`
	Content []byte      `json:",string" form:"content"`
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
type MessageBroadcastType struct {
	Content []byte
	Type    MessageType
}

var (
	BroadcastNewChannel MessageBroadcastType = MessageBroadcastType{
		Content: []byte("new channel created"),
		Type:    -1,
	}
)

type MessageType int

const (
	MessageTypeImage   MessageType = 1
	MessageTypeSticker MessageType = 2
	MessageTypeText    MessageType = 3
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
