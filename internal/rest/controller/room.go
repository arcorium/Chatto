package controller

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/util/httputil"
	"chatto/internal/util/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRoomController(roomService service.IRoomService) IController {
	return roomController{roomService: roomService}
}

type roomController struct {
	roomService service.IRoomService
}

func (r roomController) Route(router gin.IRouter, middlewares *middleware.Middleware) {
	roomRoute := router.Group("/rooms", middlewares.UserAgent.Handle(), middlewares.TokenValidation.Handle())
	roomRoute.GET("/", r.GetAllRoom)
	roomRoute.GET("/:id", r.GetRoomById)
	roomRoute.POST("/", r.CreateRoom)
	roomRoute.DELETE("/:id", r.DeleteRoomById)
}

func (r roomController) AddUserIntoRoomById(ctx *gin.Context) {
	roomId := ctx.Param("id")
	if strutil.IsEmpty(roomId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}

	var userIds dto.UserRoomInput
	if err := ctx.ShouldBind(&userIds); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	output, err := r.roomService.AddUserIntoRoom(roomId, userIds)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, output)
}

func (r roomController) RemoveUserFromRoomById(ctx *gin.Context) {
	roomId := ctx.Param("id")
	if strutil.IsEmpty(roomId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}

	var userIds dto.UserRoomInput
	if err := ctx.ShouldBind(&userIds); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	err := r.roomService.RemoveUsersOnRoom(roomId, userIds)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, nil)
}

func (r roomController) CreateRoom(ctx *gin.Context) {
	var createRoom dto.CreateRoomInput
	if err := ctx.ShouldBind(&createRoom); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	roomOutput, err := r.roomService.CreateRoom(&createRoom)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusCreated, roomOutput)
}

func (r roomController) GetAllRoom(ctx *gin.Context) {
	rm, err := r.roomService.FindRooms()
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, rm)
}

func (r roomController) GetRoomById(ctx *gin.Context) {
	roomId := ctx.Param("id")
	if strutil.IsEmpty(roomId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}

	rm, err := r.roomService.FindRoomById(roomId)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, rm)
}

func (r roomController) DeleteRoomById(ctx *gin.Context) {
	roomId := ctx.Param("id")
	if strutil.IsEmpty(roomId) {
		httputil.ErrorResponse(ctx, http.StatusInternalServerError, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}
	force := ctx.DefaultQuery("force", "0") == "1"

	err := r.roomService.DeleteRoomById(roomId, force)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, nil)
}
