package handlers

import (
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func NewAuthHandler(us *services.UserServices, as *services.ActivityServices) *AuthHandler {
	return &AuthHandler{
		Authorized:      false,
		UserService:     us,
		ActivityService: as,
	}
}

type AuthHandler struct {
	Authorized      bool
	UserService     *services.UserServices
	ActivityService *services.ActivityServices
}

func (ah *AuthHandler) successfulPost(c echo.Context) error {
	c.Response().Header().Set("HX-Push-URL", "/home")
	c.Response().Header().Set("HX-Retarget", "body")
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	return renderView(c, views.Base(views.Home(), ah.UserService.User))
}

func (ah *AuthHandler) loginPostHandler(c echo.Context) error {
	user, err := ah.UserService.CheckEmail(c.FormValue("email"))
	if err != nil {
		return c.HTML(http.StatusOK, "Bad Email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
	if err != nil {
		return c.HTML(http.StatusOK, "Wrong Password")
	}

	log.Println("Logged in", user.FirstName, user.LastName)

	ah.Authorized = true
	ah.UserService.User = user

	return ah.successfulPost(c)
}

func (ah *AuthHandler) registerPostHandler(c echo.Context) error {
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirmPassword")

	//if len(formData.Password) < 10 {
	//	formData.Message = "Password must be at least 10 characters"
	//	return renderView(c, auth.Register(formData))
	//}

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

	err := ah.UserService.CreateUser(services.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return c.HTML(http.StatusOK, "Could not create account")
	}

	ah.Authorized = true
	ah.UserService.User = services.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	return ah.successfulPost(c)
}

func (ah *AuthHandler) signoutPostHandler(c echo.Context) error {
	ah.UserService.User = services.User{}
	ah.Authorized = false

	return ah.successfulPost(c)
}
