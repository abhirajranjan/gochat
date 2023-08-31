package ports

import (
	"net/http"
)

type handlerfunc func(http.ResponseWriter, *http.Request)
type middlewarefunc func(http.Handler) http.Handler

type Handler interface {
	// handle login via google response data
	HandleGoogleAuth() handlerfunc
	// retrives recent user messages
	GetUserMessages() handlerfunc
	// Delete User
	DeleteUser() handlerfunc

	// Create New Channel
	NewChannel() handlerfunc

	// all below have channelid as parameter
	//
	// get messages from channel
	GetMessagesFromChannel() handlerfunc
	// post a message to channel
	PostMessageInChannel() handlerfunc
	// User join channel
	JoinChannel() handlerfunc
	// Delete Channel
	DeleteChannel() handlerfunc

	// injects user data in context if present else returns http.Unauthorised
	AuthMiddleware() middlewarefunc

	// upgrade connection to websocket
	HandleWS() handlerfunc
}
