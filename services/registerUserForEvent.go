package services

import (
	"github.com/jmoiron/sqlx"
)

func RegisterUserForEvent(Db *sqlx.DB, eventId int64, userId int64) error {
	if _, err := Db.Exec(`INSERT INTO eventRegistrations (event_id, user_id) VALUES (?, ?);`, eventId, userId); err != nil {
		return err
	}

	return nil
}
