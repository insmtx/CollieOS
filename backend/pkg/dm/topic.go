package dm

import (
	"strings"
	"unicode"
)

const wildcard = "*"
const unknownSegment = "unknown"

// TopicBuilder builds domain message subjects from ordered segments.
type TopicBuilder struct {
	segments []string
}

// Topic creates an empty topic builder.
func Topic() TopicBuilder {
	return TopicBuilder{}
}

// Add appends one or more literal subject segments.
func (b TopicBuilder) Add(segments ...string) TopicBuilder {
	next := b.clone()
	for _, segment := range segments {
		next.segments = append(next.segments, cleanSegment(segment))
	}
	return next
}

// Org appends organization field segments.
func (b TopicBuilder) Org(orgID string) TopicBuilder {
	return b.Add("org", orgID)
}

// Session appends session field segments.
func (b TopicBuilder) Session(sessionID string) TopicBuilder {
	return b.Add("session", sessionID)
}

// Worker appends worker field segments.
func (b TopicBuilder) Worker(workerID string) TopicBuilder {
	return b.Add("worker", workerID)
}

// Message appends the message field segment.
func (b TopicBuilder) Message() TopicBuilder {
	return b.Add("message")
}

// Stream appends the stream field segment.
func (b TopicBuilder) Stream() TopicBuilder {
	return b.Add("stream")
}

// Task appends the task field segment.
func (b TopicBuilder) Task() TopicBuilder {
	return b.Add("task")
}

// Wildcard appends a single-token wildcard segment.
func (b TopicBuilder) Wildcard() TopicBuilder {
	next := b.clone()
	next.segments = append(next.segments, wildcard)
	return next
}

// Build returns the final dot-separated subject.
func (b TopicBuilder) Build() string {
	return strings.Join(b.segments, ".")
}

func (b TopicBuilder) clone() TopicBuilder {
	segments := make([]string, len(b.segments))
	copy(segments, b.segments)
	return TopicBuilder{segments: segments}
}

func cleanSegment(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Map(func(r rune) rune {
		switch {
		case r == '.' || r == '*' || r == '>':
			return '_'
		case unicode.IsSpace(r):
			return '_'
		default:
			return r
		}
	}, value)
	if value == "" {
		return unknownSegment
	}
	return value
}
