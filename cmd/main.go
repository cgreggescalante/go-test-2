package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"nff-go-htmx/models"
	"path/filepath"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"nff-go-htmx/db"

	"io"
	"os"
	"strings"
	"time"
)

// DbName TODO: move to env file
const (
	DbName    = "gotest2.sqlite"
	SecretKey = "secret"
)

const (
	AuthSessionKey   = "auth-session"
	AuthKey          = "auth"
	UserIdKey        = "user-id"
	UserFirstNameKey = "user-first-name"
	UserLastNameKey  = "user-last-name"
)

type UserData struct {
	IsAuthorized bool
	FirstName    string
	LastName     string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		return t.templates.ExecuteTemplate(w, name, data)
	}

	authorized, ok := c.Get(AuthKey).(bool)
	if !ok {
		authorized = false
	}

	firstName, _ := c.Get(UserFirstNameKey).(string)
	lastName, _ := c.Get(UserLastNameKey).(string)

	return t.templates.ExecuteTemplate(w, "layout", struct {
		UserData
		Path     string
		PageData interface{}
	}{
		UserData: UserData{
			IsAuthorized: authorized,
			FirstName:    firstName,
			LastName:     lastName,
		},
		Path:     c.Request().URL.Path,
		PageData: data,
	})
}

func Home(c echo.Context) error {
	authorized, ok := c.Get(AuthKey).(bool)
	if !ok {
		authorized = false
	}

	return c.Render(http.StatusOK, "home", struct{ Authorized bool }{authorized})
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func CreateLoginPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		var user models.User
		checkEmailErr := db.Get(&user, "SELECT id, email, password, first_name, last_name FROM users WHERE email = ?", email)
		checkPasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if checkEmailErr != nil || checkPasswordErr != nil {
			return c.HTML(http.StatusOK, "Bad Email / Password")
		}

		sess, _ := session.Get(AuthSessionKey, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		fmt.Println("Logged in", user.FirstName, user.LastName)

		sess.Values = map[interface{}]interface{}{
			AuthKey:          true,
			UserIdKey:        user.ID,
			UserFirstNameKey: user.FirstName,
			UserLastNameKey:  user.LastName,
		}

		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}

		c.Response().Header().Set("HX-Redirect", "/")

		return nil
	}
}

func CreateRegisterPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirmPassword")

		if password != confirmPassword {
			return c.HTML(http.StatusOK, "Passwords do not match")
		}

		firstName := c.FormValue("firstName")
		if firstName == "" {
			return c.HTML(http.StatusOK, "First Name is required")
		}

		lastName := c.FormValue("lastName")
		if lastName == "" {
			return c.HTML(http.StatusOK, "Last Name is required")
		}

		email := c.FormValue("email")
		if email == "" {
			return c.HTML(http.StatusOK, "Email is required")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			return err
		}

		statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`
		result, err := db.Exec(statement, email, string(hashedPassword), firstName, lastName)
		if err != nil {
			return err
		}

		userId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		sess, _ := session.Get(AuthSessionKey, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		sess.Values = map[interface{}]interface{}{
			AuthKey:          true,
			UserIdKey:        userId,
			UserFirstNameKey: firstName,
			UserLastNameKey:  lastName,
		}

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}

		c.Response().Header().Set("HX-Redirect", "/")

		return nil
	}
}

func SignOut(c echo.Context) error {
	sess, _ := session.Get(AuthSessionKey, c)
	sess.Values = map[interface{}]interface{}{
		AuthKey:          false,
		UserIdKey:        0,
		UserFirstNameKey: "",
		UserLastNameKey:  "",
	}
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")

	return nil
}

func CreateEventListHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var events []models.Event

		if err := db.Select(&events, "SELECT id, name, description, start, end, registration_start, registration_end FROM events;"); err != nil {
			return err
		}

		return c.Render(http.StatusOK, "events", events)
	}
}

func CreateEventHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(c.Path())
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			fmt.Println(c.Param("id"))
			fmt.Printf("Failed to parse event id: %s\n", err)
			return err
		}

		var event models.Event
		if err := db.Get(&event, "SELECT * FROM events WHERE id = ?;", eventId); err != nil {
			return err
		}

		authorized, ok := c.Get(AuthKey).(bool)
		if !ok {
			authorized = false
		}

		userId, ok := c.Get(UserIdKey).(int64)
		if !ok {
			userId = 0
		}

		var count int
		if err := db.Get(&count, "SELECT COUNT(*) FROM eventRegistrations WHERE event_id = ? AND user_id = ?;", eventId, userId); err != nil {
			count = 0
		}

		data := struct {
			Event            models.Event
			Authorized       bool
			Registered       bool
			RegistrationOpen bool
		}{
			Event:            event,
			Authorized:       authorized,
			Registered:       count > 0,
			RegistrationOpen: event.RegistrationStart < time.Now().Unix() && event.RegistrationEnd > time.Now().Unix(),
		}

		return c.Render(http.StatusOK, "event", data)
	}
}

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
			if err := db.Select(&data.Activities, "SELECT u.first_name, u.last_name, activities.* FROM activities JOIN main.users u on u.id = activities.user_id ORDER BY date DESC LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize); err != nil {
				return err
			}
		} else {
			userId, ok := c.Get(UserIdKey).(int64)
			if !ok {
				return c.HTML(http.StatusOK, "Not authenticated")
			}

			firstName := c.Get(UserFirstNameKey).(string)
			lastName := c.Get(UserLastNameKey).(string)

			if err := db.Select(&data.Activities, "SELECT * FROM activities WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?", userId, pageSize, (page-1)*pageSize); err != nil {
				return err
			}

			for i := 0; i < len(data.Activities); i++ {
				data.Activities[i].FirstName = firstName
				data.Activities[i].LastName = lastName
			}
		}

		data.More = data.Activities != nil && len(data.Activities) == pageSize

		return c.Render(http.StatusOK, "recentActivitiesBlock", data)
	}
}

type AddActivityData struct {
	ActivityTypes []string
}

func AddActivity(c echo.Context) error {
	return c.Render(http.StatusOK, "addActivity", AddActivityData{ActivityTypes: models.ActivityTypes})
}

func CreateActivityPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, ok := c.Get(UserIdKey).(int64)
		if !ok {
			return c.HTML(http.StatusOK, "Not authenticated")
		}

		durations := map[string]float64{}
		foundIncluded := false

		for _, item := range models.ActivityTypes {
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

		_, err := db.NamedExec(`INSERT INTO activities (
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
					);`, activity)
		if err != nil {
			fmt.Println(err)
			return c.HTML(http.StatusOK, "Could not upload activity")
		}

		return c.HTML(http.StatusOK, "Processed")
	}
}

func CreateLeaderboardHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var entries []models.LeaderboardEntry

		if err := db.Select(&entries, "SELECT u.first_name, u.last_name, SUM(run_points + classic_roller_skiing_points + skate_roller_skiing_points + road_biking_points + mountain_biking_points + walking_points + hiking_with_packs_points + swimming_points + paddling_points + strength_training_points + aerobic_sports_points) AS points FROM activities JOIN main.users u on u.id = activities.user_id GROUP BY u.id ORDER BY points DESC;"); err != nil {
			return err
		}

		for i := 0; i < len(entries); i++ {
			if i > 0 && entries[i].Points == entries[i-1].Points {
				entries[i].Rank = entries[i-1].Rank
				continue
			}
			entries[i].Rank = i + 1
		}

		return c.Render(http.StatusOK, "leaderboard", entries)
	}
}

func CreateEventRegistrationHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}

		userId, ok := c.Get(UserIdKey).(int64)
		if !ok {
			return c.HTML(http.StatusOK, "You must be logged in to register for this event.")
		}

		var event models.Event
		if err := db.Get(&event, "SELECT * FROM events WHERE id = ?;", eventId); err != nil {
			return err
		}

		if event.RegistrationStart > time.Now().Unix() || event.RegistrationEnd < time.Now().Unix() {
			return c.HTML(http.StatusOK, "Registration is not open.")
		}

		if _, err = db.Exec(`INSERT INTO eventRegistrations (event_id, user_id) VALUES (?, ?);`, eventId, userId); err != nil {
			return err
		}

		return c.HTML(http.StatusOK, "You are registered for this event.")
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(AuthSessionKey, c)

		if auth, ok := sess.Values[AuthKey].(bool); !auth || !ok {
			c.Set(AuthKey, false)

			return next(c)
		}

		if userId, ok := sess.Values[UserIdKey].(int64); userId != 0 && ok {
			c.Set(UserIdKey, userId)
		}

		if firstName, ok := sess.Values[UserFirstNameKey].(string); firstName != "" && ok {
			c.Set(UserFirstNameKey, firstName)
		}

		if lastName, ok := sess.Values[UserLastNameKey].(string); lastName != "" && ok {
			c.Set(UserLastNameKey, lastName)
		}

		c.Set(AuthKey, true)

		return next(c)
	}
}

func CreateRoutes(e *echo.Echo, Db *sqlx.DB) {
	t := &Template{
		templates: template.New(""),
	}

	filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			_, err := t.templates.ParseFiles(path)
			if err != nil {
				fmt.Printf("Failed to parse file: %s", err)
			}
		}

		return nil
	})

	for _, template := range t.templates.Templates() {
		fmt.Printf("Template: %s\n", template.Name())
	}

	e.Renderer = t

	e.GET("/", Home)
	e.GET("/home", Home)
	e.GET("/activities", CreateActivityListHandler(Db))
	e.GET("/addActivity", AddActivity)
	e.POST("/addActivity", CreateActivityPostHandler(Db))

	e.GET("/leaderboard", CreateLeaderboardHandler(Db))

	e.GET("/events", CreateEventListHandler(Db))
	e.GET("/event/:id", CreateEventHandler(Db))
	e.POST("/event/:id/register", CreateEventRegistrationHandler(Db))

	e.GET("/login", Login)
	e.GET("/register", Register)

	e.POST("/login", CreateLoginPostHandler(Db))
	e.POST("/register", CreateRegisterPostHandler(Db))
	e.POST("/signout", SignOut)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SecretKey))))

	Db, err := db.NewStore(DbName)
	if err != nil {
		e.Logger.Fatalf("Failed to create store: %s", err)
	}

	e.Use(AuthMiddleware)

	CreateRoutes(e, Db)

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

	Db, _ := db.NewStore(DbName)

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
