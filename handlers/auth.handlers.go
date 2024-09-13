package handlers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nff-go-htmx/services"
	"nff-go-htmx/views"
)

const (
	authSessionKey   = "auth-session"
	authKey          = "auth"
	userIdKey        = "user-id"
	userFirstNameKey = "user-first-name"
	userLastNameKey  = "user-last-name"
)

func NewAuthHandler(us *services.UserServices, as *services.ActivityServices, es *services.EventServices) *AuthHandler {
	return &AuthHandler{
		UserService:     us,
		ActivityService: as,
		EventService:    es,
	}
}

type AuthHandler struct {
	UserService     *services.UserServices
	ActivityService *services.ActivityServices
	EventService    *services.EventServices
}

// AuthMiddleware adds the user's session to the context. This middleware does not block unauthorized users.
func (ah *AuthHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(authSessionKey, c)

		if auth, ok := sess.Values[authKey].(bool); !auth || !ok {
			c.Set(authKey, false)

			fmt.Println("Not authorized")

			return next(c)
		}

		if userId, ok := sess.Values[userIdKey].(int64); userId != 0 && ok {
			c.Set(userIdKey, userId)
		}

		if firstName, ok := sess.Values[userFirstNameKey].(string); firstName != "" && ok {
			c.Set(userFirstNameKey, firstName)
		}

		if lastName, ok := sess.Values[userLastNameKey].(string); lastName != "" && ok {
			c.Set(userLastNameKey, lastName)
		}

		c.Set(authKey, true)

		fmt.Println("Authorized")

		return next(c)
	}
}

func (ah *AuthHandler) successfulPost(c echo.Context, userName string, authorized bool) error {
	c.Response().Header().Set("HX-Push-URL", "/home")
	c.Response().Header().Set("HX-Retarget", "body")
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	return renderView(c, views.Base(views.Home(authorized), userName, authorized))
}

func (ah *AuthHandler) loginPostHandler(c echo.Context) error {
	user, err := ah.UserService.CheckEmail(c.FormValue("email"))
	if err != nil {
		fmt.Println("Error checking email", err)
		return c.HTML(http.StatusOK, "Bad Email / Password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password"))); err != nil {
		return c.HTML(http.StatusOK, "Bad Email / Password")
	}

	log.Println("Logged in", user.FirstName, user.LastName)

	sess, _ := session.Get(authSessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values = map[interface{}]interface{}{
		authKey:          true,
		userIdKey:        user.ID,
		userFirstNameKey: user.FirstName,
		userLastNameKey:  user.LastName,
	}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return ah.successfulPost(c, fmt.Sprintf("%s %s", user.FirstName, user.LastName), true)
}

func (ah *AuthHandler) registerPostHandler(c echo.Context) error {
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

	userId, err := ah.UserService.CreateUser(services.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return c.HTML(http.StatusOK, "Could not create account")
	}

	sess, _ := session.Get(authSessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values = map[interface{}]interface{}{
		authKey:          true,
		userIdKey:        userId,
		userFirstNameKey: firstName,
		userLastNameKey:  lastName,
	}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return ah.successfulPost(c, fmt.Sprintf("%s %s", firstName, lastName), true)
}

func (ah *AuthHandler) signoutPostHandler(c echo.Context) error {
	sess, _ := session.Get(authSessionKey, c)
	sess.Values = map[interface{}]interface{}{
		authKey:          false,
		userIdKey:        0,
		userFirstNameKey: "",
		userLastNameKey:  "",
	}
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return ah.successfulPost(c, "", false)
}
