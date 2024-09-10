package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-test-2/services"
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
