package models

type User struct {
	ID        int64
	Email     string
	Password  string
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
