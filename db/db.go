package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Store struct {
	Db *sqlx.DB
}

func NewStore(name string) (Store, error) {
	var (
		Db  *sqlx.DB
		err error
	)

	if Db != nil {
		return Store{}, nil
	}

	Db, err = sqlx.Open("sqlite3", name)
	if err != nil {
		return Store{}, fmt.Errorf("failed to connect to the database: %s", err)
	}

	if err = createUsersTable(Db); err != nil {
		return Store{}, fmt.Errorf("failed to create users table: %s", err)
	}

	if err = createActivitiesTable(Db); err != nil {
		return Store{}, fmt.Errorf("failed to create activities table: %s", err)
	}

	if err = createEventTable(Db); err != nil {
		return Store{}, fmt.Errorf("failed to create event table: %s", err)
	}

	if err = createEventRegistrationTable(Db); err != nil {
		return Store{}, fmt.Errorf("failed to create event registration table: %s", err)
	}

	log.Println("Connected to the database")

	return Store{Db}, nil
}

func createUsersTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL,
		    first_name VARCHAR(64) NOT NULL,
		    last_name VARCHAR(64) NOT NULL
		);
	`)
	return err
}

func createActivitiesTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS activities (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    user_id INTEGER,
		    date INT,
		    description VARCHAR(255),
		    run REAL,
		    run_points REAL,
			classic_roller_skiing REAL,
			classic_roller_skiing_points REAL,
			skate_roller_skiing REAL,
			skate_roller_skiing_points REAL,
			road_biking REAL,
			road_biking_points REAL,
			mountain_biking REAL,
			mountain_biking_points REAL,
			walking REAL,
			walking_points REAL,
			hiking_with_packs REAL,
			hiking_with_packs_points REAL,
			swimming REAL,
			swimming_points REAL,
			paddling REAL,
			paddling_points REAL,
			strength_training REAL,
			strength_training_points REAL,
			aerobic_sports REAL,
			aerobic_sports_points REAL,
		    FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	return err
}

func createEventTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name VARCHAR(255),
		    description TEXT,
		    start INT,
		    end INT,
		    registration_start INT,
		    registration_end INT
		)
	`)
	return err
}

func createEventRegistrationTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS eventRegistrations (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    event_id INTEGER,
		    user_id INTEGER,
		    created INT,
		    FOREIGN KEY (event_id) REFERENCES events(id),
		    FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	return err
}
