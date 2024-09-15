package handlers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"nff-go-htmx/config"
	"nff-go-htmx/models"
)

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func CreateLoginPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		var user models.User
		checkEmailErr := db.Get(&user, "SELECT id, email, password, first_name, last_name FROM users WHERE email = ?", email)
		checkPasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if checkEmailErr != nil || checkPasswordErr != nil {
			return c.HTML(http.StatusOK, "Bad Email / Password")
		}

		sess, _ := session.Get(config.AuthSessionKey, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		fmt.Println("Logged in", user.FirstName, user.LastName)

		sess.Values = map[interface{}]interface{}{
			config.AuthKey:          true,
			config.UserIdKey:        user.ID,
			config.UserFirstNameKey: user.FirstName,
			config.UserLastNameKey:  user.LastName,
		}

		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}

		c.Response().Header().Set("HX-Redirect", "/")

		return nil
	}
}
