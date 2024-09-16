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
		leaderboard, err := services.GetLeaderboard(db)
		if err != nil {
			fmt.Printf("Error in CreateLeaderboardHandler: %v\n", err)
		}

		for i := 0; i < len(leaderboard); i++ {
			leaderboard[i].Points = float64(int64(leaderboard[i].Points*100)) / 100
		}

		return c.Render(http.StatusOK, "leaderboard", leaderboard)
	}
}
