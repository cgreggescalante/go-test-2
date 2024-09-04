package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-test-2/db"
	"go-test-2/handlers"
	"go-test-2/services"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	store, err := db.NewStore("gotest2.sqlite")
	if err != nil {
		e.Logger.Fatalf("Failed to create store: %s", err)
	}

	us := services.NewUserService(services.User{}, store)
	ah := handlers.NewAuthHandler(us)

	handlers.SetRoutes(e, ah)

	e.Logger.Fatal(e.Start(":8080"))
}
