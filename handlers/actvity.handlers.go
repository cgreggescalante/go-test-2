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
	userId, ok := c.Get(userIdKey).(int64)
	if !ok {
		return c.HTML(http.StatusOK, "Not authenticated")
	}

	durations := map[string]float64{}
	foundIncluded := false

	for _, item := range services.ActivityTypes {
		duration, err := strconv.ParseFloat("0"+c.FormValue(item), 32)
		if err != nil {
			return c.HTML(http.StatusOK, fmt.Sprintf("Bad input for %s duration", item))
		}
		durations[item] = duration
		if duration > 0 {
			foundIncluded = true
		}
	}

	if !foundIncluded {
		return c.HTML(http.StatusOK, "Cannot upload activities with no values!!")
	}

	err := ah.ActivityService.CreateActivity(services.Activity{
		UserId:                    userId,
		Date:                      time.Now().Unix(),
		Description:               c.FormValue("description"),
		Run:                       durations["Run"],
		RunPoints:                 durations["Run"],
		ClassicRollerSkiing:       durations["Classic Roller Skiing"],
		ClassicRollerSkiingPoints: durations["Classic Roller Skiing"],
		SkateRollerSkiing:         durations["Skate Roller Skiing"],
		SkateRollerSkiingPoints:   durations["Skate Roller Skiing"],
		RoadBiking:                durations["Road Biking"],
		RoadBikingPoints:          durations["Road Biking"],
		MountainBiking:            durations["Mountain Biking"],
		MountainBikingPoints:      durations["Mountain Biking"],
		Walking:                   durations["Walking"],
		WalkingPoints:             durations["Walking"],
		HikingWithPacks:           durations["Hiking With Packs"],
		HikingWithPacksPoints:     durations["Hiking With Packs"],
		Swimming:                  durations["Swimming"],
		SwimmingPoints:            durations["Swimming"],
		Paddling:                  durations["Paddling"],
		PaddlingPoints:            durations["Paddling"],
		StrengthTraining:          durations["Strength Training"],
		StrengthTrainingPoints:    durations["Strength Training"],
		AerobicSports:             durations["Aerobic Sports"],
		AerobicSportsPoints:       durations["Aerobic Sports"],
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
		userId, ok := c.Get(userIdKey).(int64)
		if !ok {
			return c.HTML(http.StatusOK, "Not authenticated")
		}
		firstName, _ := c.Get(userFirstNameKey).(string)
		lastName, _ := c.Get(userLastNameKey).(string)

		activities, err = ah.ActivityService.GetRecentActivitiesByUser(services.User{ID: userId, FirstName: firstName, LastName: lastName}, page, pageSize)
	}
	if err != nil {
		fmt.Println(err)
		return c.HTML(http.StatusOK, "Could not get activities")
	}

	return renderView(c, partials.ActivityBlock(activities, page+1, pageSize, "all", len(activities) == pageSize))
}
