package models

var ActivityTypes = []string{"Run", "Classic Roller Skiing", "Skate Roller Skiing", "Road Biking", "Mountain Biking", "Walking", "Hiking With Packs", "Swimming", "Paddling", "Strength Training", "Aerobic Sports"}

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
