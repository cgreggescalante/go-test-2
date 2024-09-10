package services

import "go-test-2/db"

var ActivityTypes = [...]string{"Run", "Bike", "Ski", "Swim"}

type Activity struct {
	UserId      int
	Date        int64
	Description string
	Run         float64
	RunPoints   float64
	Bike        float64
	BikePoints  float64
	Ski         float64
	SkiPoints   float64
	Swim        float64
	SwimPoints  float64
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
	statement := `INSERT INTO activities (userId, date, description, run, runPoints, bike, bikePoints, ski, skiPoints, swim, swimPoints) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := as.ActivityStore.Db.Exec(
		statement,
		activity.UserId,
		activity.Date,
		activity.Description,
		activity.Run, activity.RunPoints,
		activity.Bike, activity.BikePoints,
		activity.Ski, activity.SkiPoints,
		activity.Swim, activity.SwimPoints,
	)
	if err != nil {
		return err
	}
	return nil
}
