package ports

import (
	"net/http"
)

type Handler interface {
	// handle login via google response data
	HandleGoogleAuth() http.Handler
	// retrives recent user messages
	GetUserMessages() http.Handler
	// Delete User
	DeleteUser() http.Handler

	// Create New Channel
	NewChannel() http.Handler

	// all below have channelid as parameter
	//
	// get messages from channel
	GetMessagesFromChannel() http.Handler
	// post a message to channel
	PostMessageInChannel() http.Handler
	// User join channel
	JoinChannel() http.Handler
	// Delete Channel
	DeleteChannel() http.Handler

	// injects user data in context if present else returns http.Unauthorised
	AuthMiddleware() func(http.Handler) http.Handler

	// upgrade connection to websocket
	HandleWS() http.Handler
}
