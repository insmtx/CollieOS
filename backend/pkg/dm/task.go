package dm

// MessagePriority describes business urgency for message routing.
type MessagePriority string

const (
	// PriorityP0 indicates an execution-critical message.
	PriorityP0 MessagePriority = "P0"
	// PriorityP1 indicates an important but non-blocking message.
	PriorityP1 MessagePriority = "P1"
	// PriorityP2 indicates a normal user-initiated message.
	PriorityP2 MessagePriority = "P2"
)

// TaskType describes the work requested from a worker.
type TaskType string

const (
	// TaskTypeAgentRun asks a worker to run an agent.
	TaskTypeAgentRun TaskType = "agent.run"
)

// InputType describes the primary shape of task input.
type InputType string

const (
	// InputTypeMessage represents normal chat input.
	InputTypeMessage InputType = "message"
	// InputTypeEvent represents an external channel event.
	InputTypeEvent InputType = "event"
	// InputTypeTaskInstruction represents a direct task instruction.
	InputTypeTaskInstruction InputType = "task_instruction"
)

// MessageRole describes who produced a chat or stream message.
type MessageRole string

const (
	// MessageRoleUser is a human or external user message.
	MessageRoleUser MessageRole = "user"
	// MessageRoleAssistant is an assistant message.
	MessageRoleAssistant MessageRole = "assistant"
	// MessageRoleSystem is a system message.
	MessageRoleSystem MessageRole = "system"
	// MessageRoleTool is a tool result message.
	MessageRoleTool MessageRole = "tool"
)

// WorkerTaskMessage is the Server -> Worker task message contract.
type WorkerTaskMessage = Envelope[WorkerTaskBody]

// WorkerTaskBody is the Server -> Worker task payload.
type WorkerTaskBody struct {
	TaskType TaskType        `json:"task_type"`
	Priority MessagePriority `json:"priority,omitempty"`

	Actor  ActorContext  `json:"actor"`
	Target TargetContext `json:"target"`
	Input  TaskInput     `json:"input"`

	Runtime RuntimeOptions `json:"runtime,omitempty"`
	Policy  TaskPolicy     `json:"policy,omitempty"`
}

// ActorContext describes the initiator of a task.
type ActorContext struct {
	UserID      string `json:"user_id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Channel     string `json:"channel,omitempty"`
	ExternalID  string `json:"external_id,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
}

// TargetContext describes the assistant, agent, and capabilities selected by the server.
type TargetContext struct {
	AssistantID string   `json:"assistant_id,omitempty"`
	AgentID     string   `json:"agent_id,omitempty"`
	Skills      []string `json:"skills,omitempty"`
	Tools       []string `json:"tools,omitempty"`
}

// TaskInput contains the normalized input consumed by the worker runtime.
type TaskInput struct {
	Type        InputType      `json:"type"`
	Text        string         `json:"text,omitempty"`
	Messages    []ChatMessage  `json:"messages,omitempty"`
	Attachments []Attachment   `json:"attachments,omitempty"`
	Event       map[string]any `json:"event,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// ChatMessage is a compact conversation message snapshot.
type ChatMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// Attachment describes an input attachment made available to a task.
type Attachment struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	URL      string `json:"url,omitempty"`
}

// RuntimeOptions controls worker runtime execution.
type RuntimeOptions struct {
	Kind    string `json:"kind,omitempty"`
	WorkDir string `json:"work_dir,omitempty"`
	MaxStep int    `json:"max_step,omitempty"`
}

// TaskPolicy carries policy knobs for a worker task.
type TaskPolicy struct {
	RequireApproval bool `json:"require_approval,omitempty"`
}
