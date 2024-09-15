package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"nff-go-htmx/config"
	"nff-go-htmx/models"
	"nff-go-htmx/services"
	"strconv"
)

type RecentActivitiesBlockData struct {
	Activities []models.ActivityWithUser
	More       bool
	Page       int
	PageSize   int
	Filter     string
}

func CreateActivityListHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
			fmt.Printf("Failed to parse page: %s", err)
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
		if err != nil {
			pageSize = 10
			fmt.Printf("Failed to parse pageSize: %s", err)
		}

		filter := c.QueryParam("user")
		if filter == "" {
			filter = "all"
		}

		var data = RecentActivitiesBlockData{
			PageSize: pageSize,
			Page:     page + 1,
			Filter:   filter,
		}

		if filter == "all" {
			activities, err := services.GetRecentActivities(db, page, pageSize)
			if err != nil {
				fmt.Printf("Error in CreateActivityListHandler: %v\n", err)
			}
			data.Activities = activities
		} else {
			userId, ok := c.Get(config.UserIdKey).(int64)
			if !ok {
				return c.HTML(http.StatusOK, "Not authenticated")
			}

			firstName := c.Get(config.UserFirstNameKey).(string)
			lastName := c.Get(config.UserLastNameKey).(string)

			activities, err := services.GetRecentActivitiesByUser(db, page, pageSize, models.User{ID: userId, FirstName: firstName, LastName: lastName})
			if err != nil {
				fmt.Printf("Error in CreateActivityListHandler: %v\n", err)
			}
			data.Activities = activities
		}

		data.More = data.Activities != nil && len(data.Activities) == pageSize

		return c.Render(http.StatusOK, "recentActivitiesBlock", data)
	}
}
