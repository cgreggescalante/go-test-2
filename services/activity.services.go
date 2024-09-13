package services

import (
	"nff-go-htmx/db"
)

var ActivityTypes = [...]string{"Run", "Classic Roller Skiing", "Skate Roller Skiing", "Road Biking", "Mountain Biking", "Walking", "Hiking With Packs", "Swimming", "Paddling", "Strength Training", "Aerobic Sports"}

type Activity struct {
	Id                        int64
	UserId                    int64 `db:"user_id"`
	Date                      int64
	Description               string
	Run                       float64
	RunPoints                 float64 `db:"run_points"`
	ClassicRollerSkiing       float64 `db:"classic_roller_skiing"`
	ClassicRollerSkiingPoints float64 `db:"classic_roller_skiing_points"`
	SkateRollerSkiing         float64 `db:"skate_roller_skiing"`
	SkateRollerSkiingPoints   float64 `db:"skate_roller_skiing_points"`
	RoadBiking                float64 `db:"road_biking"`
	RoadBikingPoints          float64 `db:"road_biking_points"`
	MountainBiking            float64 `db:"mountain_biking"`
	MountainBikingPoints      float64 `db:"mountain_biking_points"`
	Walking                   float64
	WalkingPoints             float64 `db:"walking_points"`
	HikingWithPacks           float64 `db:"hiking_with_packs"`
	HikingWithPacksPoints     float64 `db:"hiking_with_packs_points"`
	Swimming                  float64
	SwimmingPoints            float64 `db:"swimming_points"`
	Paddling                  float64
	PaddlingPoints            float64 `db:"paddling_points"`
	StrengthTraining          float64 `db:"strength_training"`
	StrengthTrainingPoints    float64 `db:"strength_training_points"`
	AerobicSports             float64 `db:"aerobic_sports"`
	AerobicSportsPoints       float64 `db:"aerobic_sports_points"`
}

type LeaderboardEntry struct {
	User
	Points float64
	Rank   int
}

type ActivityWithUser struct {
	Activity
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type ActivityServices struct {
	ActivityStore db.Store
}

func NewActivityService(activityStore db.Store) *ActivityServices {
	return &ActivityServices{
		ActivityStore: activityStore,
	}
}

func (as *ActivityServices) CreateActivity(activity Activity) error {
	statement := `INSERT INTO activities (
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
					);`

	_, err := as.ActivityStore.Db.NamedExec(statement, activity)

	if err != nil {
		return err
	}
	return nil
}

func (as *ActivityServices) GetRecentActivities(page int, pageSize int) ([]ActivityWithUser, error) {
	query := `SELECT u.first_name, u.last_name, activities.* FROM activities JOIN main.users u on u.id = activities.user_id ORDER BY date DESC LIMIT ? OFFSET ?`

	var result []ActivityWithUser

	if err := as.ActivityStore.Db.Select(&result, query, pageSize, (page-1)*pageSize); err != nil {
		return []ActivityWithUser{}, err
	}

	return result, nil
}

func (as *ActivityServices) GetRecentActivitiesByUser(user User, page int, pageSize int) ([]ActivityWithUser, error) {
	query := `SELECT * FROM activities WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?`

	var activities []ActivityWithUser
	if err := as.ActivityStore.Db.Select(&activities, query, user.ID, pageSize, (page-1)*pageSize); err != nil {
		return []ActivityWithUser{}, err
	}

	var result []ActivityWithUser
	for _, activity := range activities {
		activity.FirstName = user.FirstName
		activity.LastName = user.LastName
		result = append(result, activity)
	}

	return result, nil
}

func (as *ActivityServices) GetLeaderboard() ([]LeaderboardEntry, error) {
	query := `SELECT u.*,
       			SUM(run_points + classic_roller_skiing_points + skate_roller_skiing_points + road_biking_points + mountain_biking_points + walking_points + hiking_with_packs_points + swimming_points + paddling_points + strength_training_points + aerobic_sports_points) 
           		as points FROM activities JOIN main.users u on u.id = activities.user_id GROUP BY user_id ORDER BY points DESC`

	var data []LeaderboardEntry

	if err := as.ActivityStore.Db.Select(&data, query); err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i++ {
		if i > 0 && data[i].Points == data[i-1].Points {
			data[i].Rank = data[i-1].Rank
			continue
		}
		data[i].Rank = i + 1
	}

	return data, nil
}
