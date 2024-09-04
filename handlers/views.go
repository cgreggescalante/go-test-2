package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"go-test-2/views"
	"go-test-2/views/auth"
)

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func baseHandler(c echo.Context) error {
	return renderView(c, views.Base(views.Home()))
}

func homeHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, views.Home())
	}
	return renderView(c, views.Base(views.Home()))
}

func otherHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, views.Other())
	}
	return renderView(c, views.Base(views.Other()))
}

func loginHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, auth.Login(""))
	}
	return renderView(c, views.Base(auth.Login("")))
}

func registerHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, auth.Register(*auth.CreateRegisterFormData()))
	}
	return renderView(c, views.Base(auth.Register(*auth.CreateRegisterFormData())))
}
