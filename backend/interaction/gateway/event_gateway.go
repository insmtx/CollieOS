package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/insmtx/SingerOS/backend/interaction"
	"github.com/insmtx/SingerOS/backend/interaction/connectors/github"
	"github.com/insmtx/SingerOS/backend/interaction/eventbus"
)

func SetupRouter(r gin.IRouter, publisher eventbus.Publisher) {
	registry := interaction.NewRegistry()
	githubConnector := github.NewConnector(
	// cfg.Github,
	// publisher,
	)
	registry.Register(githubConnector)

	registry.RegisterRoutes(r)
}
