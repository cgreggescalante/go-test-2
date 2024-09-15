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

func CreateEventListHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		events, err := services.GetEvents(db)

		if err != nil {
			fmt.Printf("Error in CreateEventListHandler: %v\n", err)
		}

		return c.Render(http.StatusOK, "events", events)
	}
}

func CreateEventHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.HTML(http.StatusOK, "Could not find event.")
		}

		event, err := services.GetEvent(db, eventId)
		if err != nil {
			return c.HTML(http.StatusOK, "Could not find event.")
		}

		authorized, ok := c.Get(config.AuthKey).(bool)
		if !ok {
			authorized = false
		}

		userId, ok := c.Get(config.UserIdKey).(int64)
		if !ok {
			userId = 0
		}

		registered, err := services.CheckUserRegisteredForEvent(db, eventId, userId)
		if err != nil {
			registered = false
		}

		data := struct {
			Event            models.Event
			Authorized       bool
			Registered       bool
			RegistrationOpen bool
		}{
			Event:            event,
			Authorized:       authorized,
			Registered:       registered,
			RegistrationOpen: event.RegistrationStart < time.Now().Unix() && event.RegistrationEnd > time.Now().Unix(),
		}

		return c.Render(http.StatusOK, "event", data)
	}
}

func CreateEventRegistrationHandler(db *sqlx.DB) echo.HandlerFunc {
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

		if err := services.RegisterUserForEvent(db, eventId, userId); err != nil {
			fmt.Printf("Error in CreateEventRegistrationHandler: %v\n", err)
			return c.HTML(http.StatusOK, "Could not register.")
		}

		return c.HTML(http.StatusOK, "You are registered for this event.")
	}
}
