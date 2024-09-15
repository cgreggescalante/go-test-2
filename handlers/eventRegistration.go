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

type EventRegistrationData struct {
	Event     models.Event
	Divisions []string
}

func EventRegistration(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}

		registrationOpen, err := services.CheckRegistrationOpen(db, eventId)
		if err != nil {
			fmt.Printf("Error in EventRegistration: %v\n", err)
			return c.HTML(http.StatusOK, "Could not register.")
		}
		if !registrationOpen {
			return c.HTML(http.StatusOK, "Registration is not open for this event.")
		}

		event, err := services.GetEvent(db, eventId)
		if err != nil {
			return c.HTML(http.StatusOK, "Could not find event.")
		}

		data := EventRegistrationData{
			Event:     event,
			Divisions: models.Divisions,
		}

		return c.Render(http.StatusOK, "eventRegistration", data)
	}
}

func CreateEventRegistrationHandler(db *sqlx.DB) echo.HandlerFunc {
	eventHandler := CreateEventHandler(db)

	return func(c echo.Context) error {
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}

		userId, ok := c.Get(config.UserIdKey).(int64)
		if !ok {
			return c.HTML(http.StatusOK, "You must be logged in to register for this event.")
		}

		registrationOpen, err := services.CheckRegistrationOpen(db, eventId)
		if err != nil {
			fmt.Printf("Error in CreateEventRegistrationHandler: %v\n", err)
			return c.HTML(http.StatusOK, "Could not register.")
		}
		if !registrationOpen {
			return c.HTML(http.StatusOK, "Registration is not open for this event.")
		}

		goalString := c.FormValue("goal")
		if goalString == "" {
			return c.HTML(http.StatusOK, "You must provide a goal.")
		}
		goal, err := strconv.ParseFloat(goalString, 64)
		if err != nil {
			return c.HTML(http.StatusOK, "Goal must be a number.")
		}

		division := c.FormValue("division")
		if division == "" {
			return c.HTML(http.StatusOK, "You must provide a division.")
		}

		registration := models.EventRegistration{
			EventId:  eventId,
			UserId:   userId,
			Division: c.FormValue("division"),
			Goal:     goal,
			Created:  time.Now().Unix(),
		}

		if err := services.RegisterUserForEvent(db, registration); err != nil {
			fmt.Printf("Error in CreateEventRegistrationHandler: %v\n", err)
			return c.HTML(http.StatusOK, "Could not register.")
		}

		c.Response().Header().Set("HX-Push-URL", fmt.Sprintf("/event/%d", eventId))
		c.Response().Header().Set("HX-Retarget", "main")

		return eventHandler(c)
	}
}
