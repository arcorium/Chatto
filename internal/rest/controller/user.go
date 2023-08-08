package controller

import (
	"log"
	"net/http"

	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/util/httputil"
	"chatto/internal/util/strutil"

	"chatto/internal/model"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/util"

	"github.com/gin-gonic/gin"
)

func NewUserController(service service.IUserService) IController {
	return &userController{service: service}
}

type userController struct {
	service service.IUserService
}

func (u userController) Route(router gin.IRouter, middlewares *middleware.Middleware) {
	userRoute := router.Group("/users", middlewares.UserAgent, middlewares.TokenValidation)
	userRoute.GET("/:id", u.GetUserById)
	userRoute.PUT("/:id", u.UpdateUser)
	userRoute.DELETE("/:id", u.RemoveUser)

	userRoute.Use(middlewares.AdminPrivilege)
	userRoute.GET("/", u.GetUsers)
	userRoute.POST("/", u.CreateUser)
}

func (u userController) CreateUser(ctx *gin.Context) {
	var user dto.CreateUserInput
	if err := ctx.BindJSON(&user); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	err := u.service.CreateUser(&user)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusCreated, nil)
}

func (u userController) GetUsers(ctx *gin.Context) {
	users, err := u.service.GetUsers()
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, users)
}

func (u userController) GetUserById(ctx *gin.Context) {
	userId := ctx.Param("id")
	if strutil.IsEmpty(userId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}

	// Special Case
	if userId == "me" {
		accessClaims, err := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)
		// TODO: Delete error checking, it is not necessary
		if err != nil {
			log.Println(err)
		}
		userId = accessClaims.UserId
	}

	user, err := u.service.FindUserById(userId)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, user)
}

func (u userController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	if strutil.IsEmpty(userId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}
	if userId == "me" {
		accessClaims, _ := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)
		userId = accessClaims.UserId
	}

	var user dto.UpdateUserInput
	if err := ctx.BindJSON(&user); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	err := u.service.UpdateUserById(userId, &user)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, nil)
}

func (u userController) RemoveUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	if strutil.IsEmpty(userId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}
	err := u.service.RemoveUserById(userId)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, nil)
}
