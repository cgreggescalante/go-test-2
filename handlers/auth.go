package handlers

import (
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(u services.User) error
	CheckEmail(email string) (services.User, error)
}

func NewAuthHandler(us *services.UserServices) *AuthHandler {
	return &AuthHandler{
		UserService: us,
	}
}

type AuthHandler struct {
	UserService *services.UserServices
}

func (ah *AuthHandler) loginPostHandler(c echo.Context) error {
	user, err := ah.UserService.CheckEmail(c.FormValue("email"))
	if err != nil {
		return renderView(c, auth.Login("Bad Email"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
	if err != nil {
		return renderView(c, auth.Login("Wrong Password"))
	}

	return renderView(c, auth.Login("Success"))
}

func (ah *AuthHandler) registerPostHandler(c echo.Context) error {
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirmPassword")

	formData := auth.RegisterFormData{
		FirstName:       c.FormValue("firstName"),
		LastName:        c.FormValue("lastName"),
		Email:           c.FormValue("email"),
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	if len(password) < 10 {
		formData.Message = "Password must be at least 10 characters"
		return renderView(c, auth.Register(formData))
	}

	if password != confirmPassword {
		formData.Message = "Passwords do not match"
		return renderView(c, auth.Register(formData))
	}

	formData.Message = "Success"

	return renderView(c, auth.Register(formData))
}
