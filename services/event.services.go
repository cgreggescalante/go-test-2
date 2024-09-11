package services

import (
	"go-test-2/db"
)

type Event struct {
	Id                int64
	Name              string
	Description       string
	Start             int64
	End               int64
	RegistrationStart int64 `db:"registration_start"`
	RegistrationEnd   int64 `db:"registration_end"`
}

type EventRegistration struct {
	Id      int64
	EventId int64
	UserId  int64
	Created int64
}

type EventServices struct {
	EventStore db.Store
}

func NewEventService(eventStore db.Store) *EventServices {
	return &EventServices{
		EventStore: eventStore,
	}
}

func (es *EventServices) RegisterUser(eventId int64, userId int64) error {
	statement := `INSERT INTO eventRegistrations (event_id, user_id) VALUES (?, ?);`

	_, err := es.EventStore.Db.Exec(statement, eventId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (es *EventServices) GetEvent(id int64) (Event, error) {
	var event Event

	if err := es.EventStore.Db.Get(&event, `SELECT * FROM events WHERE id = ?;`, id); err != nil {
		return Event{}, err
	}

	return event, nil
}

func (es *EventServices) GetEvents() ([]Event, error) {
	rows, err := es.EventStore.Db.Query(`SELECT id, name, description, start, end, registration_start, registration_end FROM events;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Start, &event.End, &event.RegistrationStart, &event.RegistrationEnd)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (es *EventServices) CheckRegistration(eventId int64, userId int64) (bool, error) {
	var count int
	if err := es.EventStore.Db.Get(&count, `SELECT COUNT(*) FROM eventRegistrations WHERE event_id = ? AND user_id = ?;`, eventId, userId); err != nil {
		return false, err
	}

	return count > 0, nil
}
