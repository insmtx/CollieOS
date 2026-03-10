package types

import "gorm.io/gorm"

type Event struct {
	gorm.Model

	MessageID string
	TraceID   string
	Source    string
	Type      string
	Action    string

	Actor  string
	Target string

	Payload map[string]any

	Timestamp int64
}
