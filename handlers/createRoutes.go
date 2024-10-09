package handlers

import (
	"bytes"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
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

// MINIFY HTML TEMPLATES FOR MASSIVE SAVINGS
func compileTemplates() (*template.Template, error) {
	paths := make([]string, 0)

	filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			paths = append(paths, path)
		}
		return nil
	})

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
	htmlMinifier := &html.Minifier{TemplateDelims: html.GoTemplateDelims}

	var tmpl *template.Template
	for _, path := range paths {
		name := filepath.Base(path)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		r := bytes.NewBufferString(string(b))
		w := &bytes.Buffer{}
		err = htmlMinifier.Minify(m, w, r, nil)
		if err != nil {
			return nil, err
		}

		tmpl.Parse(w.String())
	}
	return tmpl, nil
}

func CreateRoutes(e *echo.Echo, Db *sqlx.DB) {
	t := &Template{
		templates: template.Must(compileTemplates()),
	}

	e.Renderer = t

	e.GET("/", Home)
	e.GET("/home", Home)
	e.GET("/uploads", CreateUploadListHandler(Db))
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
