package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"nff-go-htmx/services"
)

func CreateLeaderboardHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		entries, err := services.GetLeaderboard(db)
		if err != nil {
			fmt.Printf("Error in CreateLeaderboardHandler: %v\n", err)
		}

		return c.Render(http.StatusOK, "leaderboard", entries)
	}
}
