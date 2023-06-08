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

func NewUserController(service service.IUserService) *UserController {
	return &UserController{service: service}
}

type UserController struct {
	service service.IUserService
}

func (u *UserController) CreateUser(ctx *gin.Context) {
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

func (u *UserController) GetUsers(ctx *gin.Context) {
	users, err := u.service.FindUsers()
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, users)
}

func (u *UserController) GetUserById(ctx *gin.Context) {
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

func (u *UserController) UpdateUser(ctx *gin.Context) {
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

func (u *UserController) RemoveUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := u.service.RemoveUserById(id)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, id)
}
