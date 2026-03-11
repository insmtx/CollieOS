package interaction

import (
	"github.com/gin-gonic/gin"
	"github.com/insmtx/SingerOS/backend/interaction/connectors"
)

type Registry struct {
	connectors map[string]connectors.Connector
}

func NewRegistry() *Registry {

	return &Registry{
		connectors: map[string]connectors.Connector{},
	}
}

func (r *Registry) Register(c connectors.Connector) {
	r.connectors[c.ChannelCode()] = c
}

func (r *Registry) RegisterRoutes(
	router gin.IRouter,
) {

	for _, c := range r.connectors {
		c.RegisterRoutes(router)
	}
}
