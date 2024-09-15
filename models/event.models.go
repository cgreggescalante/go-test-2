package models

type Event struct {
	Id                int64
	Name              string
	Description       string
	Start             int64
	End               int64
	RegistrationStart int64 `db:"registration_start"`
	RegistrationEnd   int64 `db:"registration_end"`
}
