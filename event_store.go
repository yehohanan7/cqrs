package cqrs

type EventStore interface {
	GetAllEvents() []Event
	GetEvents(aggregateId string) []Event
	GetEvent(id string) Event
	SaveEvents(aggregateId string, events []Event) error
}
