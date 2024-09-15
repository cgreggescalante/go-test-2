package services

import "github.com/jmoiron/sqlx"

func CheckUserRegisteredForEvent(db *sqlx.DB, eventId int64, userId int64) (bool, error) {
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM eventRegistrations WHERE event_id = ? AND user_id = ?;", eventId, userId); err != nil {
		return false, err
	}

	return count > 0, nil
}
