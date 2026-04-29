package dm

import (
	"strings"
	"unicode"
)

const wildcard = "*"
const unknownSegment = "unknown"

// TopicBuilder 基于有序片段构造领域消息 topic。
type TopicBuilder struct {
	segments []string
}

// Topic 创建一个空的 topic builder。
func Topic() TopicBuilder {
	return TopicBuilder{}
}

// Add 追加一个或多个普通 topic 片段。
func (b TopicBuilder) Add(segments ...string) TopicBuilder {
	next := b.clone()
	for _, segment := range segments {
		next.segments = append(next.segments, cleanSegment(segment))
	}
	return next
}

// Org 追加组织字段片段。
func (b TopicBuilder) Org(orgID string) TopicBuilder {
	return b.Add("org", orgID)
}

// Session 追加会话字段片段。
func (b TopicBuilder) Session(sessionID string) TopicBuilder {
	return b.Add("session", sessionID)
}

// Worker 追加 Worker 字段片段。
func (b TopicBuilder) Worker(workerID string) TopicBuilder {
	return b.Add("worker", workerID)
}

// Message 追加 message 字段片段。
func (b TopicBuilder) Message() TopicBuilder {
	return b.Add("message")
}

// Stream 追加 stream 字段片段。
func (b TopicBuilder) Stream() TopicBuilder {
	return b.Add("stream")
}

// Task 追加 task 字段片段。
func (b TopicBuilder) Task() TopicBuilder {
	return b.Add("task")
}

// Wildcard 追加单层通配符片段。
func (b TopicBuilder) Wildcard() TopicBuilder {
	next := b.clone()
	next.segments = append(next.segments, wildcard)
	return next
}

// Build 返回使用点号连接后的最终 topic。
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
