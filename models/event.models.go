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

var Divisions = []string{"Upperclassmen", "Underclassmen", "Middle School", "Staff & VIPs", "Parents", "Alumni"}

type EventRegistration struct {
	Id       int64
	EventId  int64 `db:"event_id"`
	UserId   int64 `db:"user_id"`
	Division string
	Goal     float64
	Created  int64
}
