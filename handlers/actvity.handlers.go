package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-test-2/services"
	"go-test-2/views/partials"
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

	fmt.Println(ah.UserService.User.ID)

	err := ah.ActivityService.CreateActivity(services.Activity{
		UserId:                    ah.UserService.User.ID,
		Date:                      time.Now().Unix(),
		Description:               c.FormValue("description"),
		Run:                       durations[0],
		RunPoints:                 durations[0],
		ClassicRollerSkiing:       durations[1],
		ClassicRollerSkiingPoints: durations[1],
		SkateRollerSkiing:         durations[2],
		SkateRollerSkiingPoints:   durations[2],
		RoadBiking:                durations[3],
		RoadBikingPoints:          durations[3],
		MountainBiking:            durations[4],
		MountainBikingPoints:      durations[4],
		Walking:                   durations[5],
		WalkingPoints:             durations[5],
		HikingWithPacks:           durations[6],
		HikingWithPacksPoints:     durations[6],
		Swimming:                  durations[7],
		SwimmingPoints:            durations[7],
		Paddling:                  durations[8],
		PaddlingPoints:            durations[8],
		StrengthTraining:          durations[9],
		StrengthTrainingPoints:    durations[9],
		AerobicSports:             durations[10],
		AerobicSportsPoints:       durations[10],
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
