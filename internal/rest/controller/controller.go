package controller

import (
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetUsers(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	RemoveUser(ctx *gin.Context)
}

type IAuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
	LogoutAllDevice(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type IChatController interface {
	Chat(ctx *gin.Context)
}
