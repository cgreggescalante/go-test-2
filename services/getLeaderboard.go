package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func GetLeaderboard(Db *sqlx.DB) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry

	if err := Db.Select(&entries, "SELECT u.first_name, u.last_name, SUM(run_points + classic_roller_skiing_points + skate_roller_skiing_points + road_biking_points + mountain_biking_points + walking_points + hiking_with_packs_points + swimming_points + paddling_points + strength_training_points + aerobic_sports_points) AS points FROM activities JOIN main.users u on u.id = activities.user_id GROUP BY u.id ORDER BY points DESC;"); err != nil {
		fmt.Printf("Error in GetLeaderboard: %v\n", err)
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
