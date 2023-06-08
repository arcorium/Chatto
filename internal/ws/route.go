package ws

import (
	"chatto/internal/config"
	"chatto/internal/rest/middleware"
	"chatto/internal/ws/handler"

	"github.com/gin-gonic/gin"
)

func NewWebsocketRoute(websocketHandler handler.WebsocketHandler) Route {
	return Route{handler: websocketHandler}
}

type Route struct {
	handler handler.WebsocketHandler
}

func (r Route) RegisterRoute(cfg *config.AppConfig, router gin.IRouter) {
	authMiddleware := middleware.TokenValidationMiddleware{Config: &middleware.TokenValidationConfig{
		SecretKeyFunc: cfg.JWTKeyFunc,
		TokenType:     "Bearer",
		SigningType:   "HS512",
	}}

	router.GET("/chat", authMiddleware.Handle(), r.handler.ServeWebsocket)
}
