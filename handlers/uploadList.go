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

type ActivityData struct {
	Name     string
	Duration float64
	Points   float64
}

type RecentActivitiesActivityData struct {
	Activities    []ActivityData
	FirstName     string
	LastName      string
	Description   string
	DateFormatted string
}

type RecentUploadsBlockData struct {
	Uploads  []RecentActivitiesActivityData
	More     bool
	Page     int
	PageSize int
	Filter   string
}

func CreateUploadListHandler(db *sqlx.DB) echo.HandlerFunc {
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

		var data = RecentUploadsBlockData{
			PageSize: pageSize,
			Page:     page + 1,
			Filter:   filter,
		}

		if filter == "all" {
			uploads, err := services.GetRecentUploads(db, page, pageSize)
			if err != nil {
				fmt.Printf("Error in CreateUploadListHandler: %v\n", err)
			}
			data.Uploads = UploadsToActivityData(uploads)
		} else {
			userId, ok := c.Get(config.UserIdKey).(int64)
			if !ok {
				return c.HTML(http.StatusOK, "Not authenticated")
			}

			firstName := c.Get(config.UserFirstNameKey).(string)
			lastName := c.Get(config.UserLastNameKey).(string)

			uploads, err := services.GetRecentUploadsByUser(db, page, pageSize, models.User{ID: userId, FirstName: firstName, LastName: lastName})
			if err != nil {
				fmt.Printf("Error in CreateUploadListHandler: %v\n", err)
			}
			data.Uploads = UploadsToActivityData(uploads)
		}

		data.More = data.Uploads != nil && len(data.Uploads) == pageSize

		return c.Render(http.StatusOK, "recentUploadsBlock", data)
	}
}

func UploadsToActivityData(uploads []models.UploadWithUser) []RecentActivitiesActivityData {
	var result []RecentActivitiesActivityData

	for _, upload := range uploads {
		var activities []ActivityData
		for _, activityType := range models.ActivityTypes {
			duration := upload.GetDuration(activityType)
			if duration > 0 {
				activities = append(activities, ActivityData{Name: activityType, Duration: duration, Points: upload.GetPoints(activityType)})
			}
		}
		result = append(result, RecentActivitiesActivityData{
			Activities:    activities,
			FirstName:     upload.FirstName,
			LastName:      upload.LastName,
			DateFormatted: upload.DateObj.Format("2006-01-02"),
		})
	}

	return result
}
