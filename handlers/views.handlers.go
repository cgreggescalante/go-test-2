package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"go-test-2/views"
	"go-test-2/views/auth"
	"go-test-2/views/events"
	"strconv"
)

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (ah *AuthHandler) fullPageRender(c echo.Context, cmp templ.Component) error {
	return renderView(c, views.Base(cmp, ah.UserService.User))
}

func (ah *AuthHandler) baseHandler(c echo.Context) error {
	return ah.fullPageRender(c, views.Home(ah.Authorized))
}

func (ah *AuthHandler) rerenderBody(c echo.Context, cmp templ.Component) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, cmp)
	}
	return renderView(c, views.Base(cmp, ah.UserService.User))
}

func (ah *AuthHandler) homeHandler(c echo.Context) error {
	return ah.rerenderBody(c, views.Home(ah.Authorized))
}

func (ah *AuthHandler) addActivityHandler(c echo.Context) error {
	return ah.rerenderBody(c, views.AddActivity())
}

func (ah *AuthHandler) loginHandler(c echo.Context) error {
	return ah.rerenderBody(c, auth.Login())
}

func (ah *AuthHandler) registerHandler(c echo.Context) error {
	return ah.rerenderBody(c, auth.Register())
}

func (ah *AuthHandler) leaderboardHandler(c echo.Context) error {
	return ah.rerenderBody(c, views.Leaderboard(ah.UserService.User))
}

func (ah *AuthHandler) getLeaderboardHandler(c echo.Context) error {
	data, err := ah.ActivityService.GetLeaderboard()
	if err != nil {
		return err
	}

	return renderView(c, views.LeaderboardBody(data))
}

func (ah *AuthHandler) eventsHandler(c echo.Context) error {
	data, err := ah.EventService.GetEvents()
	if err != nil {
		return err
	}

	return ah.rerenderBody(c, events.EventList(data))
}

func (ah *AuthHandler) eventHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	data, err := ah.EventService.GetEvent(id)
	if err != nil {
		return err
	}

	registered := false

	if ah.Authorized {
		registered, err = ah.EventService.CheckRegistration(id, ah.UserService.User.ID)
		if err != nil {
			return err
		}
	}

	return ah.rerenderBody(c, events.Event(data, ah.Authorized, registered))
}

func (ah *AuthHandler) registerEventHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	err = ah.EventService.RegisterUser(id, ah.UserService.User.ID)
	if err != nil {
		return err
	}

	return renderView(c, events.RegistrationStatus(id, true, true))
}
