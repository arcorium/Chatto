package route

import (
	"server_client_chat/internal/config"

	"github.com/gin-gonic/gin"
)

type IRoute interface {
	V1Handle(router gin.IRouter, cfg *config.AppConfig)
}

func V1Route(router gin.IRouter, cfg *config.AppConfig, routes ...IRoute) {
	v1 := router.Group("/api/v1")
	for _, x := range routes {
		x.V1Handle(v1, cfg)
	}
}
