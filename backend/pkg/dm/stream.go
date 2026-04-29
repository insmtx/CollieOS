package dm

// StreamEventType identifies one event in a worker run stream.
type StreamEventType string

const (
	// StreamEventRunStarted indicates a run has started.
	StreamEventRunStarted StreamEventType = "run.started"
	// StreamEventMessageDelta contains assistant output delta text.
	StreamEventMessageDelta StreamEventType = "message.delta"
	// StreamEventToolCallStarted indicates a tool call has started.
	StreamEventToolCallStarted StreamEventType = "tool_call.started"
	// StreamEventToolCallDelta contains streamed tool call arguments.
	StreamEventToolCallDelta StreamEventType = "tool_call.delta"
	// StreamEventToolCallFinished indicates a tool call has finished.
	StreamEventToolCallFinished StreamEventType = "tool_call.finished"
	// StreamEventMessageCompleted contains the final assistant message.
	StreamEventMessageCompleted StreamEventType = "message.completed"
	// StreamEventRunCompleted indicates a run completed successfully.
	StreamEventRunCompleted StreamEventType = "run.completed"
	// StreamEventRunFailed indicates a run failed.
	StreamEventRunFailed StreamEventType = "run.failed"
)

// MessageStreamMessage is the Worker -> Server -> UI stream message contract.
type MessageStreamMessage = Envelope[StreamBody]

// StreamBody is one event in a Worker -> Server -> UI stream.
type StreamBody struct {
	Seq     int64           `json:"seq"`
	Event   StreamEventType `json:"event"`
	Payload StreamPayload   `json:"payload"`

	Usage *UsagePayload `json:"usage,omitempty"`
	Error *StreamError  `json:"error,omitempty"`
}

// StreamPayload contains the event-specific stream content.
type StreamPayload struct {
	Role       MessageRole      `json:"role,omitempty"`
	Content    string           `json:"content,omitempty"`
	ToolCall   *ToolCallEvent   `json:"tool_call,omitempty"`
	ToolResult *ToolResultEvent `json:"tool_result,omitempty"`
}

// ToolCallEvent describes a tool invocation emitted during a stream.
type ToolCallEvent struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

// ToolResultEvent describes a tool result emitted during a stream.
type ToolResultEvent struct {
	ToolCallID string         `json:"tool_call_id"`
	Name       string         `json:"name,omitempty"`
	Result     map[string]any `json:"result,omitempty"`
}

// UsagePayload describes model usage when available.
type UsagePayload struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	TotalTokens  int `json:"total_tokens,omitempty"`
}

// StreamError describes a terminal or recoverable stream error.
type StreamError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
