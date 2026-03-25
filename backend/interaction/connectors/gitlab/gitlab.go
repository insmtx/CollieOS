package gitlab

import (
	"github.com/insmtx/SingerOS/backend/config"
	"github.com/insmtx/SingerOS/backend/interaction/connectors"
	"github.com/insmtx/SingerOS/backend/interaction/eventbus"
)

var _ connectors.Connector = (*GitlabConnector)(nil)

type GitlabConnector struct {
	config    config.GitlabAppConfig
	publisher eventbus.Publisher
}

func NewConnector(cfg config.GitlabAppConfig, publisher eventbus.Publisher) *GitlabConnector {
	return &GitlabConnector{
		config:    cfg,
		publisher: publisher,
	}
}

func (c *GitlabConnector) ChannelCode() string {
	return "gitlab"
}
