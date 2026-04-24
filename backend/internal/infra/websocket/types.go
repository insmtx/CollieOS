package websocket

import "time"

const ChannelCodeValue = "websocket"

// Message represents possible message types from clients
type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	ID      string                 `json:"id,omitempty"`
}

// ServerMessage represents possible message types sent to clients
type ServerMessage struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	ID        string                 `json:"id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// MessageDestination specifies who should receive the message
type MessageDestination struct {
	ClientID string `json:"client_id"`
	TaskID   string `json:"task_id"`
}
