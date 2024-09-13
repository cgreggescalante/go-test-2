package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"nff-go-htmx/views"
	"nff-go-htmx/views/auth"
	"nff-go-htmx/views/events"
	"strconv"
)

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (ah *AuthHandler) rerenderBody(c echo.Context, cmp templ.Component) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return renderView(c, cmp)
	}

	authorized, ok := c.Get(authKey).(bool)
	if !ok {
		authorized = false
	}

	userFirstName, _ := c.Get(userFirstNameKey).(string)
	userLastName, _ := c.Get(userLastNameKey).(string)

	return renderView(c, views.Base(cmp, fmt.Sprintf("%s %s", userFirstName, userLastName), authorized))
}

func (ah *AuthHandler) homeHandler(c echo.Context) error {
	authorized, ok := c.Get(authKey).(bool)
	if !ok {
		authorized = false
	}

	return ah.rerenderBody(c, views.Home(authorized))
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
	data, err := ah.ActivityService.GetLeaderboard()
	if err != nil {
		return err
	}

	return ah.rerenderBody(c, views.Leaderboard(data))
}

func (ah *AuthHandler) eventsHandler(c echo.Context) error {
	data, err := ah.EventService.GetEvents()
	if err != nil {
		return err
	}

	return ah.rerenderBody(c, events.EventList(data))
}

func (ah *AuthHandler) eventHandler(c echo.Context) error {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	userId, ok := c.Get(userIdKey).(int64)
	if !ok {
		userId = 0
	}

	authorized, ok := c.Get(authKey).(bool)
	if !ok {
		authorized = false
	}

	data, err := ah.EventService.GetEvent(eventId)
	if err != nil {
		return err
	}

	registered := false

	if authorized {
		registered, err = ah.EventService.CheckRegistration(eventId, userId)
		if err != nil {
			return err
		}
	}

	return ah.rerenderBody(c, events.Event(data, authorized, registered))
}

func (ah *AuthHandler) registerEventHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	userId, ok := c.Get(userIdKey).(int64)
	if !ok {
		userId = 0
	}

	err = ah.EventService.RegisterUser(id, userId)
	if err != nil {
		return err
	}

	return renderView(c, events.RegistrationStatus(id, true, true, true))
}
