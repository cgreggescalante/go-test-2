package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", ah.baseHandler)
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
	e.GET("/getLeaderboard", ah.getLeaderboardHandler)

	e.GET("/events", ah.eventsHandler)
	e.GET("/event/:id", ah.eventHandler)

	e.POST("/event/:id/register", ah.registerEventHandler)
}
