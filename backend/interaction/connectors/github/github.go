package github

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v78/github"
	"github.com/insmtx/SingerOS/backend/config"
	"github.com/insmtx/SingerOS/backend/interaction/connectors"
	"github.com/insmtx/SingerOS/backend/interaction/eventbus"
	"go.uber.org/zap"
)

var _ connectors.Connector = (*GitHubConnector)(nil)

type GitHubConnector struct {
	config config.GithubAppConfig

	client *github.Client

	publisher eventbus.Publisher

	logger *zap.Logger
}

func (GitHubConnector) ChannelCode() string {
	return "github"
}

func (c *GitHubConnector) RegisterRoutes(r gin.IRouter) {
	r.POST("/github/webhook", c.HandleWebhook)
}

func NewConnector() *GitHubConnector {
	return nil
}
