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
	"time"
)

func AddActivity(c echo.Context) error {
	return c.Render(http.StatusOK, "addActivity", models.ActivityTypes)
}

func CreateActivityPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, ok := c.Get(config.UserIdKey).(int64)
		if !ok {
			return c.HTML(http.StatusOK, "Not authenticated")
		}

		durations := map[string]float64{}
		foundIncluded := false

		for _, item := range models.ActivityTypes {
			duration, err := strconv.ParseFloat("0"+c.FormValue(item), 32)
			if err != nil {
				fmt.Printf("Error in CreateActivityPostHandler: %v\n", err)
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

		activity := models.Activity{
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
		}

		if err := services.AddActivity(db, activity); err != nil {
			fmt.Printf("Error in CreateActivityPostHandler: %v\n", err)
			return c.HTML(http.StatusOK, "Could not upload activity")
		}

		return c.HTML(http.StatusOK, "Processed")
	}
}
