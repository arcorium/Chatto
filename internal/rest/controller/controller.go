package controller

import (
	"chatto/internal/rest/middleware"
	"github.com/gin-gonic/gin"
)

type IController interface {
	Route(router gin.IRouter, middleware *middleware.Middleware)
}
