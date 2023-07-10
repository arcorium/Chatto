package controller

import (
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRoomController(roomService service.IRoomService) IController {
	return roomController{roomService: roomService}
}

type roomController struct {
	roomService service.IRoomService
}

func (r roomController) Route(router gin.IRouter, middlewares *middleware.Middleware) {
	//TODO implement me
	roomRoute := router.Group("/rooms", middlewares.TokenValidation.Handle())
	roomRoute.GET("/", r.GetRoom)
	roomRoute.POST("/", r.CreateRoom)
	roomRoute.DELETE("/:id", r.DeleteRoom)
}

func (r roomController) CreateRoom(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r roomController) GetRoom(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r roomController) DeleteRoom(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
