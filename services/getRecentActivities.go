package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
	"time"
)

func GetRecentUploads(Db *sqlx.DB, page int, pageSize int) ([]models.UploadWithUser, error) {
	var activities []models.UploadWithUser

	if err := Db.Select(&activities, "SELECT u.first_name, u.last_name, a.* FROM activities a JOIN users u ON a.user_id = u.id ORDER BY date DESC LIMIT ? OFFSET ?;", pageSize, (page-1)*pageSize); err != nil {
		return []models.UploadWithUser{}, err
	}

	for i, activity := range activities {
		activities[i].DateObj = time.Unix(activity.Date, 0)
	}

	return activities, nil
}

func GetRecentUploadsByUser(Db *sqlx.DB, page int, pageSize int, user models.User) ([]models.UploadWithUser, error) {
	var activities []models.UploadWithUser
	if err := Db.Select(&activities, "SELECT * FROM activities WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?;", user.ID, pageSize, (page-1)*pageSize); err != nil {
		return nil, err
	}

	for i := range activities {
		activities[i].FirstName = user.FirstName
		activities[i].LastName = user.LastName
		activities[i].DateObj = time.Unix(activities[i].Date, 0)
	}

	return activities, nil
}
