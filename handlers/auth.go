package handlers

import (
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views"
	"go-test-2/views/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func NewAuthHandler(us *services.UserServices) *AuthHandler {
	return &AuthHandler{
		UserService: us,
	}
}

type AuthHandler struct {
	UserService *services.UserServices
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
		return renderView(c, auth.PartialLogin("Bad Email"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
	if err != nil {
		return renderView(c, auth.PartialLogin("Wrong Password"))
	}

	log.Println("Logged in", user.FirstName, user.LastName)

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
		return renderView(c, auth.PartialRegister("Passwords do not match"))
	}

	firstName := c.FormValue("firstName")
	if firstName == "" {
		return renderView(c, auth.PartialRegister("First Name is required"))
	}

	lastName := c.FormValue("lastName")
	if lastName == "" {
		return renderView(c, auth.PartialRegister("Last Name is required"))
	}

	email := c.FormValue("email")
	if email == "" {
		return renderView(c, auth.PartialRegister("Email is required"))
	}

	err := ah.UserService.CreateUser(services.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return renderView(c, auth.PartialRegister("Could not create account"))
	}

	ah.UserService.User = services.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	return ah.successfulPost(c)
}

func (ah *AuthHandler) signoutPostHandler(c echo.Context) error {
	ah.UserService.User = services.User{}

	return ah.successfulPost(c)
}
