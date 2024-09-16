package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func GetLeaderboard(Db *sqlx.DB) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry

	if err := Db.Select(&entries, "SELECT u.first_name, u.last_name, SUM(points) AS points FROM activities JOIN main.users u on u.id = activities.user_id GROUP BY u.id ORDER BY SUM(points) DESC;"); err != nil {
		return []models.LeaderboardEntry{}, err
	}

	for i := 0; i < len(entries); i++ {
		if i > 0 && entries[i].Points == entries[i-1].Points {
			entries[i].Rank = entries[i-1].Rank
			continue
		}
		entries[i].Rank = i + 1
	}

	return entries, nil
}

func GetEventLeaderboard(Db *sqlx.DB, eventId int64) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry

	if err := Db.Select(&entries, "SELECT u.id, u.first_name, u.last_name, SUM(a.points) as points FROM users u JOIN eventRegistrations er ON u.id = er.user_id JOIN activities a on u.id = a.user_id JOIN events e on er.event_id = e.id WHERE er.event_id = ? AND a.date BETWEEN e.start AND e.end GROUP BY u.id, u.first_name, u.last_name ORDER BY SUM(a.points) DESC;", eventId); err != nil {
		return []models.LeaderboardEntry{}, err
	}

	for i := 0; i < len(entries); i++ {
		if i > 0 && entries[i].Points == entries[i-1].Points {
			entries[i].Rank = entries[i-1].Rank
			continue
		}
		entries[i].Rank = i + 1
	}

	return entries, nil
}
