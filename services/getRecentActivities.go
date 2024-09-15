package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func GetRecentActivities(Db *sqlx.DB, page int, pageSize int) ([]models.ActivityWithUser, error) {
	var activities []models.ActivityWithUser

	if err := Db.Select(&activities, "SELECT u.first_name, u.last_name, a.* FROM activities a JOIN users u ON a.user_id = u.id ORDER BY date DESC LIMIT ? OFFSET ?;", pageSize, (page-1)*pageSize); err != nil {
		return []models.ActivityWithUser{}, err
	}

	return activities, nil
}

func GetRecentActivitiesByUser(Db *sqlx.DB, page int, pageSize int, user models.User) ([]models.ActivityWithUser, error) {
	var activities []models.ActivityWithUser
	if err := Db.Select(&activities, "SELECT * FROM activities WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?;", user.ID, pageSize, (page-1)*pageSize); err != nil {
		return nil, err
	}

	for i := range activities {
		activities[i].FirstName = user.FirstName
		activities[i].LastName = user.LastName
	}

	return activities, nil
}
