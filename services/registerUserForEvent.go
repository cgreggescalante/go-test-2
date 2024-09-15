package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func RegisterUserForEvent(Db *sqlx.DB, eventId int64, userId int64) error {
	if _, err := Db.Exec(`INSERT INTO eventRegistrations (event_id, user_id) VALUES (?, ?);`, eventId, userId); err != nil {
		fmt.Println("Error in RegisterUserForEvent: ", err)
		return err
	}

	return nil
}
