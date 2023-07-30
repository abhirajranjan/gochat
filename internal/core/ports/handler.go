package ports

import "github.com/gin-gonic/gin"

type Handler interface {
	// handle login via google response data
	HandleGoogleAuth(*gin.Context)
	// retrives recent user messages
	GetUserMessages(*gin.Context)

	// get messages from channel
	GetMessagesFromChannel(*gin.Context)
	// post a message to channel
	PostMessageInChannel(*gin.Context)

	// injects user data in context if present else returns http.Unauthorised
	AuthMiddleware(*gin.Context)

	// upgrade connection to websocket
	HandleWS(*gin.Context)
}
