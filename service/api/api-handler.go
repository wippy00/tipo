package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//login
	rt.router.POST("/login", rt.logIn)

	// User routes
	rt.router.GET("/users", rt.getUsers)
	rt.router.PUT("/setUserName", rt.updateUserName)
	rt.router.PUT("/setUserPhoto", rt.updateUserPhoto)

	rt.router.GET("/user/:id", rt.getUser)

	// Conversation routes
	rt.router.GET("/conversations/:id", rt.getConversation)
	rt.router.PUT("/conversations/:id/setGroupName", rt.updateConversationName)
	rt.router.PUT("/conversations/:id/setGroupPhoto", rt.updateConversationPhoto)
	rt.router.POST("/conversations/:conversation_id/addToGroup/:user_id", rt.addUserToConversation)
	rt.router.DELETE("/conversations/:conversation_id/leave", rt.removeUserFromConversation)

	rt.router.GET("/conversations", rt.getConversationOfUser)

	// Message routes
	rt.router.GET("/conversations/:id/messages", rt.getMessagesOfConversation)
	rt.router.POST("/conversations/:conversation_id/messages", rt.sendMessage)
	rt.router.POST("/conversations/:conversation_id/messages/:message_id/forward", rt.forwardMessage)
	rt.router.DELETE("/conversations/:conversation_id/messages/:message_id", rt.deleteMessage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
