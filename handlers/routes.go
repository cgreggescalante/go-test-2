package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", ah.baseHandler)
	e.GET("/home", ah.homeHandler)
	e.GET("/addActivity", ah.addActivityHandler)

	e.GET("/login", ah.loginHandler)
	e.POST("/login", ah.loginPostHandler)

	e.GET("/register", ah.registerHandler)
	e.POST("/register", ah.registerPostHandler)

	e.POST("/signout", ah.signoutPostHandler)
}
