// Package dm defines domain messaging contracts shared by SingerOS services.
package dm

import "time"

// MessageType describes the top-level domain message kind.
type MessageType string

const (
	// MessageTypeWorkerTask identifies a Server -> Worker task message.
	MessageTypeWorkerTask MessageType = "worker.task"
	// MessageTypeStream identifies a Worker -> Server -> UI stream message.
	MessageTypeStream MessageType = "message.stream"
)

// TraceContext carries correlation IDs across UI, server, worker, and runtime.
type TraceContext struct {
	TraceID   string `json:"trace_id"`
	RequestID string `json:"request_id,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
	RunID     string `json:"run_id,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
}

// RouteContext carries delivery and tenancy information for one message.
type RouteContext struct {
	OrgID     string `json:"org_id"`
	SessionID string `json:"session_id,omitempty"`
	WorkerID  string `json:"worker_id,omitempty"`
}

// Envelope is the common domain message wrapper used on MQ topics.
type Envelope[T any] struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	CreatedAt time.Time   `json:"created_at"`

	Trace TraceContext `json:"trace"`
	Route RouteContext `json:"route"`

	Body     T              `json:"body"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
