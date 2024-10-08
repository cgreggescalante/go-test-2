package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"nff-go-htmx/config"
	"nff-go-htmx/handlers"
	"nff-go-htmx/models"

	"nff-go-htmx/db"

	"io"
	"os"
	"strings"
	"time"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.AuthSessionKey, c)

		if auth, ok := sess.Values[config.AuthKey].(bool); !auth || !ok {
			c.Set(config.AuthKey, false)

			return next(c)
		}

		if userId, ok := sess.Values[config.UserIdKey].(int64); userId != 0 && ok {
			c.Set(config.UserIdKey, userId)
		}

		if firstName, ok := sess.Values[config.UserFirstNameKey].(string); firstName != "" && ok {
			c.Set(config.UserFirstNameKey, firstName)
		}

		if lastName, ok := sess.Values[config.UserLastNameKey].(string); lastName != "" && ok {
			c.Set(config.UserLastNameKey, lastName)
		}

		c.Set(config.AuthKey, true)

		return next(c)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SecretKey))))
	e.Use(middleware.Gzip())

	Db, err := db.NewStore(config.DbName)
	if err != nil {
		e.Logger.Fatalf("Failed to create store: %s", err)
	}

	e.Use(AuthMiddleware)

	handlers.CreateRoutes(e, Db)

	e.Logger.Fatal(e.Start(":8080"))
}

func loadFromCSV() {
	file, err := os.Open("uploads(1).json")
	if err != nil {
		fmt.Printf("Failed to open file: %s", err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	var data []map[string]interface{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Printf("Failed to unmarshal JSON: %s", err)
	}

	Db, _ := db.NewStore(config.DbName)

	groupedByUser := make(map[string][]map[string]any)

	for _, d := range data {
		groupedByUser[d["userId"].(string)] = append(groupedByUser[d["userId"].(string)], d)
	}

	for user, activities := range groupedByUser {
		fmt.Printf("User: %s\n", user)

		parts := strings.Split(activities[0]["userDisplayName"].(string), " ")

		firstName := parts[0]
		lastName := strings.Join(parts[1:], " ")

		statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`
		result, err := Db.Exec(statement, user, "", firstName, lastName)
		if err != nil {
			fmt.Printf("Failed to insert user: %s", err)
		}
		id, _ := result.LastInsertId()

		statement = `INSERT INTO activities (
		               user_id, date, description,
		               run, run_points,
		               classic_roller_skiing, classic_roller_skiing_points,
		               skate_roller_skiing, skate_roller_skiing_points,
		               road_biking, road_biking_points,
		               mountain_biking, mountain_biking_points,
		               walking, walking_points,
		               hiking_with_packs, hiking_with_packs_points,
		               swimming, swimming_points,
		               paddling, paddling_points,
		               strength_training, strength_training_points,
		               aerobic_sports, aerobic_sports_points
					) VALUES (
					  	:user_id, :date, :description,
						:run, :run_points,
						:classic_roller_skiing, :classic_roller_skiing_points,
						:skate_roller_skiing, :skate_roller_skiing_points,
						:road_biking, :road_biking_points,
						:mountain_biking, :mountain_biking_points,
						:walking, :walking_points,
						:hiking_with_packs, :hiking_with_packs_points,
						:swimming, :swimming_points,
						:paddling, :paddling_points,
						:strength_training, :strength_training_points,
						:aerobic_sports, :aerobic_sports_points
					);`

		for _, d := range activities {
			durations := d["activities"].(map[string]interface{})
			activityPoints := d["activityPoints"].(map[string]interface{})

			date, _ := time.Parse(time.RFC3339Nano, d["date"].(string))
			run, _ := getIfExists("Run", durations)
			runPoints, _ := getIfExists("Run", activityPoints)
			classicRollerSkiing, _ := getIfExists("Classic Roller Skiing", durations)
			classicRollerSkiingPoints, _ := getIfExists("Classic Roller Skiing", activityPoints)
			skateRollerSkiing, _ := getIfExists("Skate Roller Skiing", durations)
			skateRollerSkiingPoints, _ := getIfExists("Skate Roller Skiing", activityPoints)
			roadBiking, _ := getIfExists("Road Biking", durations)
			roadBikingPoints, _ := getIfExists("Road Biking", activityPoints)
			mountainBiking, _ := getIfExists("Mountain Biking", durations)
			mountainBikingPoints, _ := getIfExists("Mountain Biking", activityPoints)
			walking, _ := getIfExists("Walking", durations)
			walkingPoints, _ := getIfExists("Walking", activityPoints)
			hikingWithPacks, _ := getIfExists("Hiking With Packs", durations)
			hikingWithPacksPoints, _ := getIfExists("Hiking With Packs", activityPoints)
			swimming, _ := getIfExists("Swimming", durations)
			swimmingPoints, _ := getIfExists("Swimming", activityPoints)
			paddling, _ := getIfExists("Paddling", durations)
			paddlingPoints, _ := getIfExists("Paddling", activityPoints)
			strengthTraining, _ := getIfExists("Strength Training", durations)
			strengthTrainingPoints, _ := getIfExists("Strength Training", activityPoints)
			aerobicSports, _ := getIfExists("Aerobic Sports", durations)
			aerobicSportsPoints, _ := getIfExists("Aerobic Sports", activityPoints)

			_, err := Db.NamedExec(statement, &models.Activity{
				UserId:                    id,
				Date:                      date.Unix(),
				Description:               d["description"].(string),
				Run:                       run,
				RunPoints:                 runPoints,
				ClassicRollerSkiing:       classicRollerSkiing,
				ClassicRollerSkiingPoints: classicRollerSkiingPoints,
				SkateRollerSkiing:         skateRollerSkiing,
				SkateRollerSkiingPoints:   skateRollerSkiingPoints,
				RoadBiking:                roadBiking,
				RoadBikingPoints:          roadBikingPoints,
				MountainBiking:            mountainBiking,
				MountainBikingPoints:      mountainBikingPoints,
				Walking:                   walking,
				WalkingPoints:             walkingPoints,
				HikingWithPacks:           hikingWithPacks,
				HikingWithPacksPoints:     hikingWithPacksPoints,
				Swimming:                  swimming,
				SwimmingPoints:            swimmingPoints,
				Paddling:                  paddling,
				PaddlingPoints:            paddlingPoints,
				StrengthTraining:          strengthTraining,
				StrengthTrainingPoints:    strengthTrainingPoints,
				AerobicSports:             aerobicSports,
				AerobicSportsPoints:       aerobicSportsPoints,
			})
			if err != nil {
				fmt.Printf("Failed to insert activity: %s", err)
			}
		}
	}
}

func getIfExists(key string, d map[string]interface{}) (float64, error) {
	if v, ok := d[key]; ok && v != nil {
		return d[key].(float64), nil
	}
	return 0, nil
}
