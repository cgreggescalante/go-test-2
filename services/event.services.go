package services

import (
	"go-test-2/db"
)

type Event struct {
	Id                int
	Name              string
	Description       string
	Start             int64
	End               int64
	RegistrationStart int64
	RegistrationEnd   int64
}

type EventRegistration struct {
	Id      int
	EventId int
	UserId  int
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

func (es *EventServices) RegisterUser(eventId int, userId int) error {
	statement := `INSERT INTO eventRegistrations (event_id, user_id) VALUES (?, ?);`

	_, err := es.EventStore.Db.Exec(
		statement,
		eventId,
		userId,
	)
	if err != nil {
		return err
	}

	return nil
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
