package services

import (
	"go-test-2/db"
)

var ActivityTypes = [...]string{"Run", "Bike", "Ski", "Swim"}

type Activity struct {
	Id          int
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

type LeaderboardEntry struct {
	UserName string
	Points   float64
	Rank     int
}

type ActivityWithUser struct {
	Activity Activity
	User     User
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

func (as *ActivityServices) GetRecentActivities(page int, pageSize int) ([]ActivityWithUser, error) {
	query := `SELECT u.firstName, u.lastName, activities.* FROM activities JOIN main.users u on u.id = activities.userId ORDER BY date DESC LIMIT ? OFFSET ?`

	statement, err := as.ActivityStore.Db.Prepare(query)
	if err != nil {
		return []ActivityWithUser{}, err
	}

	rows, err := statement.Query(pageSize, (page-1)*pageSize)
	if err != nil {
		return []ActivityWithUser{}, err
	}

	var result []ActivityWithUser

	for rows.Next() {
		var activity Activity
		var user User
		err = rows.Scan(
			&user.FirstName, &user.LastName,
			&activity.Id,
			&activity.UserId,
			&activity.Date,
			&activity.Description,
			&activity.Run, &activity.RunPoints,
			&activity.Bike, &activity.BikePoints,
			&activity.Ski, &activity.SkiPoints,
			&activity.Swim, &activity.SwimPoints,
		)
		if err != nil {
			return []ActivityWithUser{}, err
		}

		result = append(result, ActivityWithUser{Activity: activity, User: user})
	}

	return result, nil
}

func (as *ActivityServices) GetRecentActivitiesByUser(user User, page int, pageSize int) ([]ActivityWithUser, error) {
	query := `SELECT * FROM activities WHERE userId = ? ORDER BY date DESC LIMIT ? OFFSET ?`

	statement, err := as.ActivityStore.Db.Prepare(query)
	if err != nil {
		return []ActivityWithUser{}, err
	}

	rows, err := statement.Query(user.ID, pageSize, (page-1)*pageSize)
	if err != nil {
		return []ActivityWithUser{}, err
	}

	var result []ActivityWithUser
	for rows.Next() {
		var activity Activity
		err = rows.Scan(
			&activity.Id,
			&activity.UserId,
			&activity.Date,
			&activity.Description,
			&activity.Run, &activity.RunPoints,
			&activity.Bike, &activity.BikePoints,
			&activity.Ski, &activity.SkiPoints,
			&activity.Swim, &activity.SwimPoints,
		)
		if err != nil {
			return []ActivityWithUser{}, err
		}
		result = append(result, ActivityWithUser{Activity: activity, User: user})
	}

	return result, nil
}

func (as *ActivityServices) GetLeaderboard() ([]LeaderboardEntry, error) {
	query := `SELECT u.firstName || ' ' || u.lastName, SUM(runPoints + bikePoints + skiPoints + swimPoints) as points FROM activities JOIN main.users u on u.id = activities.userId GROUP BY userId ORDER BY points DESC`

	statement, err := as.ActivityStore.Db.Prepare(query)
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	rows, err := statement.Query()
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	var result []LeaderboardEntry
	rank := 1
	for rows.Next() {
		var entry LeaderboardEntry
		err = rows.Scan(&entry.UserName, &entry.Points)
		if err != nil {
			return []LeaderboardEntry{}, err
		}
		entry.Rank = rank
		rank++
		result = append(result, entry)
	}

	return result, nil
}
