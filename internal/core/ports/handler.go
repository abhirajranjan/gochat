package ports

import "github.com/gin-gonic/gin"

type Handler interface {
	// handle login via google response data
	HandleGoogleAuth(*gin.Context)
	// retrives recent user messages
	GetUserMessages(*gin.Context)
	// Delete User
	DeleteUser(*gin.Context)

	// all below have channelid as parameter
	//
	// get messages from channel
	GetMessagesFromChannel(*gin.Context)
	// post a message to channel
	PostMessageInChannel(*gin.Context)
	// User join channel
	JoinChannel(*gin.Context)
	// Create New Channel
	NewChannel(*gin.Context)
	// Delete Channel
	DeleteChannel(*gin.Context)

	// injects user data in context if present else returns http.Unauthorised
	AuthMiddleware(*gin.Context)

	// upgrade connection to websocket
	HandleWS(*gin.Context)
}
