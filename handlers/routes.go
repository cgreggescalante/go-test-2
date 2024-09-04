package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", baseHandler)
	e.GET("/home", homeHandler)
	e.GET("/other", otherHandler)

	e.GET("/login", loginHandler)
	e.POST("/login", ah.loginPostHandler)

	e.GET("/register", registerHandler)
	e.POST("/register", ah.registerPostHandler)
}
