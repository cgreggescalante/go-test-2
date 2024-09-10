package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views/partials"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (ah *AuthHandler) addActivityPostHandler(c echo.Context) error {
	if !ah.Authorized {
		return c.HTML(http.StatusOK, "Not authenticated")
	}

	durations := make([]float64, 0, len(services.ActivityTypes))
	foundIncluded := false

	for _, item := range services.ActivityTypes {
		duration, err := strconv.ParseFloat("0"+c.FormValue(item), 32)
		log.Println(duration, err, c.FormValue(item), item)
		if err != nil {
			return c.HTML(http.StatusOK, fmt.Sprintf("Bad input for %s duration", item))
		}
		durations = append(durations, duration)
		if duration > 0 {
			foundIncluded = true
		}
	}

	if !foundIncluded {
		return c.HTML(http.StatusOK, "Cannot upload activities with no values!!")
	}

	err := ah.ActivityService.CreateActivity(services.Activity{
		UserId:      ah.UserService.User.ID,
		Date:        time.Now().Unix(),
		Description: c.FormValue("description"),
		Run:         durations[0],
		RunPoints:   durations[0],
		Bike:        durations[1],
		BikePoints:  durations[1],
		Ski:         durations[2],
		SkiPoints:   durations[2],
		Swim:        durations[3],
		SwimPoints:  durations[3],
	})
	if err != nil {
		fmt.Println(err)
		return c.HTML(http.StatusOK, "Could not upload activity")
	}

	return c.HTML(http.StatusOK, "Processed")
}

func (ah *AuthHandler) getActivityHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.HTML(http.StatusOK, "Bad value for page")
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		return c.HTML(http.StatusOK, "Bad value for pageSize")
	}
	user := c.QueryParam("user")
	if !(user == "all" || user == "current") {
		fmt.Printf("Bad value for user: %s\n", user)
		user = "all"
	}

	var activities []services.ActivityWithUser

	if user == "all" {
		activities, err = ah.ActivityService.GetRecentActivities(page, pageSize)
	} else {
		activities, err = ah.ActivityService.GetRecentActivitiesByUser(ah.UserService.User, page, pageSize)
	}
	if err != nil {
		fmt.Println(err)
		return c.HTML(http.StatusOK, "Could not get activities")
	}

	return renderView(c, partials.ActivityBlock(activities, page+1, pageSize, "all", len(activities) == pageSize))
}
