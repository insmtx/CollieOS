package connectors

import "github.com/gin-gonic/gin"

type Connector interface {
	ChannelCode() string

	RegisterRoutes(r gin.IRouter)
}
