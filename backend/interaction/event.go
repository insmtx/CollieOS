package interaction

import "time"

// Event
type Event struct {
	EventID string
	TraceID string

	Channel string

	EventType string

	Actor string

	Repository string

	Context map[string]interface{}

	Payload interface{}

	CreatedAt time.Time
}
