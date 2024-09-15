package services

import (
	"github.com/jmoiron/sqlx"
	"nff-go-htmx/models"
)

func AddActivity(db *sqlx.DB, activity models.Activity) error {
	_, err := db.NamedExec(`INSERT INTO activities (
		               user_id, date, description,
		               run, run_points,
		               classic_roller_skiing, classic_roller_skiing_points,
		               skate_roller_skiing, skate_roller_skiing_points,
		               road_biking, road_biking_points,
		               mountain_biking, mountain_biking_points,
		               walking, walking_points,
		               hiking_with_packs, hiking_with_packs_points,
		               swimming, swimming_points,
		               paddling, paddling_points,
		               strength_training, strength_training_points,
		               aerobic_sports, aerobic_sports_points
					) VALUES (
					  	:user_id, :date, :description,
						:run, :run_points,
						:classic_roller_skiing, :classic_roller_skiing_points,
						:skate_roller_skiing, :skate_roller_skiing_points,
						:road_biking, :road_biking_points,
						:mountain_biking, :mountain_biking_points,
						:walking, :walking_points,
						:hiking_with_packs, :hiking_with_packs_points,
						:swimming, :swimming_points,
						:paddling, :paddling_points,
						:strength_training, :strength_training_points,
						:aerobic_sports, :aerobic_sports_points
					);`, activity)

	return err
}
