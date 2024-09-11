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
