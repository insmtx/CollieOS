package github

import (
	"github.com/google/go-github/v78/github"
	"github.com/insmtx/SingerOS/backend/interaction"
)

// convertEvent
func (c *GitHubConnector) convertEvent(
	eventType string,
	event any,
) *interaction.Event {
	switch eventType {
	case "issue_comment":
		return c.convertIssueComment(event.(*github.IssueCommentEvent))
	default:
		return nil
	}
}

func (c *GitHubConnector) convertIssueComment(
	event *github.IssueCommentEvent,
) *interaction.Event {

	return &interaction.Event{
		Channel:    c.ChannelCode(),
		EventType:  EventTypeIssueComment,
		Actor:      event.GetComment().GetUser().GetLogin(),
		Repository: event.GetRepo().GetFullName(),
		Payload: map[string]interface{}{
			"issue_number": event.GetIssue().GetNumber(),
			"comment":      event.GetComment().GetBody(),
			"comment_id":   event.GetComment().GetID(),
		},
	}
}
