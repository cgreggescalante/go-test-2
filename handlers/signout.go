package handlers

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"nff-go-htmx/config"
)

func SignOut(c echo.Context) error {
	sess, _ := session.Get(config.AuthSessionKey, c)
	sess.Values = map[interface{}]interface{}{
		config.AuthKey:          false,
		config.UserIdKey:        0,
		config.UserFirstNameKey: "",
		config.UserLastNameKey:  "",
	}
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")

	return nil
}
