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

	// Conversation routes

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
