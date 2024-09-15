package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
	"time"
)

func CheckRegistrationOpen(Db *sqlx.DB, eventId int64) (bool, error) {
	var event models.Event
	if err := Db.Get(&event, "SELECT registration_start, registration_end FROM events WHERE id = ?;", eventId); err != nil {
		fmt.Printf("Error in CheckRegistrationOpen: %v\n", err)
		return false, err
	}

	return event.RegistrationStart < time.Now().Unix() && event.RegistrationEnd > time.Now().Unix(), nil
}
