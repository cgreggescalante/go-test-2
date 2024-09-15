package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func GetEvents(db *sqlx.DB) ([]models.Event, error) {
	var events []models.Event

	if err := db.Select(&events, "SELECT * FROM events;"); err != nil {
		return []models.Event{}, err
	}

	return events, nil
}

func GetEvent(db *sqlx.DB, eventId int64) (models.Event, error) {
	var event models.Event

	if err := db.Get(&event, "SELECT * FROM events WHERE id = ?;", eventId); err != nil {
		return models.Event{}, err
	}

	return event, nil
}
