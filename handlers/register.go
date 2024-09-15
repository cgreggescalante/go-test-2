package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"nff-go-htmx/config"
)

func Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func CreateRegisterPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirmPassword")

		if password != confirmPassword {
			return c.HTML(http.StatusOK, "Passwords do not match")
		}

		firstName := c.FormValue("firstName")
		if firstName == "" {
			return c.HTML(http.StatusOK, "First Name is required")
		}

		lastName := c.FormValue("lastName")
		if lastName == "" {
			return c.HTML(http.StatusOK, "Last Name is required")
		}

		email := c.FormValue("email")
		if email == "" {
			return c.HTML(http.StatusOK, "Email is required")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			return err
		}

		statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`
		result, err := db.Exec(statement, email, string(hashedPassword), firstName, lastName)
		if err != nil {
			return err
		}

		userId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		sess, _ := session.Get(config.AuthSessionKey, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		sess.Values = map[interface{}]interface{}{
			config.AuthKey:          true,
			config.UserIdKey:        userId,
			config.UserFirstNameKey: firstName,
			config.UserLastNameKey:  lastName,
		}

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}

		c.Response().Header().Set("HX-Redirect", "/")

		return nil
	}
}
