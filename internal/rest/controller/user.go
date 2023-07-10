package controller

import (
	"log"
	"net/http"

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
	userRoute := router.Group("/users", middlewares.TokenValidation.Handle())
	userRoute.GET("/", u.GetUsers)
	userRoute.GET("/:id", u.GetUserById)
	userRoute.PUT("/:id", u.UpdateUser)
	userRoute.DELETE("/:id", u.RemoveUser)
}

func (u userController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := u.service.CreateUser(&user)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusCreated, nil)
}

func (u userController) GetUsers(ctx *gin.Context) {
	users, err := u.service.FindUsers()
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, users)
}

func (u userController) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	// Special Case
	if id == "me" {
		accessClaims, err := util.GetContextValue[model.AccessTokenClaims](middleware.KEY_JWT_CLAIMS, ctx)
		if err != nil {
			log.Println(err)
		}
		id = accessClaims.UserId
	}

	user, err := u.service.FindUserById(id)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, user)
}

func (u userController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := u.service.UpdateUserById(id, &user)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, user)
}

func (u userController) RemoveUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := u.service.RemoveUserById(id)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, id)
}
