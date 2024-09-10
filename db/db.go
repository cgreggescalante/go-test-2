package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Store struct {
	Db *sql.DB
}

func NewStore(name string) (Store, error) {
	var (
		Db  *sql.DB
		err error
	)

	if Db != nil {
		return Store{}, nil
	}

	Db, err = sql.Open("sqlite3", name)
	if err != nil {
		return Store{}, fmt.Errorf("failed to connect to the database: %s", err)
	}

	_, err = Db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL,
		    firstName VARCHAR(64) NOT NULL,
		    lastName VARCHAR(64) NOT NULL
		);
	`)
	if err != nil {
		return Store{}, err
	}

	_, err = Db.Exec(`
		CREATE TABLE IF NOT EXISTS activities (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    userId INTEGER,
		    date INT,
		    description VARCHAR(255),
		    run REAL,
		    runPoints REAL,
			bike REAL,
			bikePoints REAL,
			ski REAL,
			skiPoints REAL,
			swim REAL,
			swimPoints REAL,
		    FOREIGN KEY (userId) REFERENCES users(id)
		)
	`)

	log.Println("Connected to the database")

	return Store{Db}, nil
}
