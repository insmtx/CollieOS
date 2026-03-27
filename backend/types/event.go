package types

import "gorm.io/gorm"

type Event struct {
	gorm.Model

	MessageID string
	TraceID   string
	Source    string
	Type      string // 建议使用 types.EventType 定义的常量值
	Action    string // 建议使用 types.EventAction 定义的常量值

	Actor  string
	Target string

	Payload map[string]interface{} `gorm:"type:jsonb"`

	Timestamp int64
}
