package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nff-go-htmx/config"
)

func Home(c echo.Context) error {
	authorized, ok := c.Get(config.AuthKey).(bool)
	if !ok {
		authorized = false
	}

	return c.Render(http.StatusOK, "home", struct{ Authorized bool }{authorized})
}
