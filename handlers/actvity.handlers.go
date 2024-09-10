package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views"
	"log"
	"strconv"
	"time"
)

func (ah *AuthHandler) addActivityPostHandler(c echo.Context) error {
	if !ah.Authorized {
		return renderView(c, views.PartialAddActivity("Not authenticated"))
	}

	durations := make([]float64, 0, len(services.ActivityTypes))
	foundIncluded := false

	for _, item := range services.ActivityTypes {
		duration, err := strconv.ParseFloat("0"+c.FormValue(item), 32)
		log.Println(duration, err, c.FormValue(item), item)
		if err != nil {
			return renderView(c, views.PartialAddActivity(fmt.Sprintf("Bad input for %s duration", item)))
		}
		durations = append(durations, duration)
		if duration > 0 {
			foundIncluded = true
		}
	}

	if !foundIncluded {
		return renderView(c, views.PartialAddActivity("Cannot upload activities with no values!!"))
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
		return renderView(c, views.PartialAddActivity("Could not upload activity"))
	}

	return renderView(c, views.PartialAddActivity("Processed"))
}
