package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"nff-go-htmx/config"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	templates *template.Template
}

type UserData struct {
	IsAuthorized bool
	FirstName    string
	LastName     string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		return t.templates.ExecuteTemplate(w, name, data)
	}

	authorized, ok := c.Get(config.AuthKey).(bool)
	if !ok {
		authorized = false
	}

	firstName, _ := c.Get(config.UserFirstNameKey).(string)
	lastName, _ := c.Get(config.UserLastNameKey).(string)

	return t.templates.ExecuteTemplate(w, "layout", struct {
		UserData
		Path     string
		PageData interface{}
	}{
		UserData: UserData{
			IsAuthorized: authorized,
			FirstName:    firstName,
			LastName:     lastName,
		},
		Path:     c.Request().URL.Path,
		PageData: data,
	})
}

func CreateRoutes(e *echo.Echo, Db *sqlx.DB) {
	t := &Template{
		templates: template.New(""),
	}

	filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			_, err := t.templates.ParseFiles(path)
			if err != nil {
				fmt.Printf("Failed to parse file: %s", err)
			}
		}

		return nil
	})

	for _, templ := range t.templates.Templates() {
		fmt.Printf("Template: %s\n", templ.Name())
	}

	e.Renderer = t

	e.GET("/", Home)
	e.GET("/home", Home)
	e.GET("/activities", CreateActivityListHandler(Db))
	e.GET("/addActivity", AddActivity)
	e.POST("/addActivity", CreateActivityPostHandler(Db))

	e.GET("/leaderboard", CreateLeaderboardHandler(Db))

	e.GET("/events", CreateEventListHandler(Db))
	e.GET("/event/:id", CreateEventHandler(Db))
	e.GET("/event/:id/register", EventRegistration(Db))
	e.POST("/event/:id/register", CreateEventRegistrationHandler(Db))

	e.GET("/login", Login)
	e.POST("/login", CreateLoginPostHandler(Db))

	e.GET("/register", Register)
	e.POST("/register", CreateRegisterPostHandler(Db))

	e.POST("/signout", SignOut)
}
