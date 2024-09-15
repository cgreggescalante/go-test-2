package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func RegisterUserForEvent(Db *sqlx.DB, registration models.EventRegistration) error {
	if _, err := Db.NamedExec(`INSERT INTO eventRegistrations (event_id, user_id, division, goal, created) VALUES (:event_id, :user_id, :division, :goal, :created);`, registration); err != nil {
		return err
	}

	return nil
}
