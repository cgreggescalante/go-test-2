package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", ah.homeHandler)
	e.GET("/home", ah.homeHandler)

	e.GET("/addActivity", ah.addActivityHandler)
	e.POST("/addActivity", ah.addActivityPostHandler)

	e.GET("/login", ah.loginHandler)
	e.POST("/login", ah.loginPostHandler)

	e.GET("/register", ah.registerHandler)
	e.POST("/register", ah.registerPostHandler)

	e.POST("/signout", ah.signoutPostHandler)

	e.GET("/activities", ah.getActivityHandler)

	e.GET("/leaderboard", ah.leaderboardHandler)

	eventGroup := e.Group("/event")

	eventGroup.GET("", ah.eventsHandler)
	eventGroup.GET("/:id", ah.eventHandler)
	eventGroup.POST("/:id/register", ah.registerEventHandler)
}
