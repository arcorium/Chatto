package controller

import (
	"chatto/internal/constant"
	"chatto/internal/model/common"

	"log"
	"net/http"

	"chatto/internal/model"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/util"

	"github.com/gin-gonic/gin"
)

func NewAuthController(service service.IAuthService) IController {
	return &authController{service: service}
}

type authController struct {
	service service.IAuthService
}

func (a authController) Route(router gin.IRouter, middlewares *middleware.Middleware) {
	authRoute := router.Group("/auth")
	authRoute.POST("/login", middlewares.UserAgent.Handle(), a.Login)
	authRoute.POST("/register", middlewares.UserAgent.Handle(), a.Register)
	authRoute.POST("/refresh", a.RefreshToken)

	authRoute.Use(middlewares.TokenValidation.Handle())
	authRoute.POST("/logout", a.Logout)
	authRoute.POST("/logout/all", a.LogoutAllDevice)
}

func (a authController) Login(ctx *gin.Context) {
	var signInInput model.SignInInput
	if err := ctx.BindJSON(&signInInput); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	systemInfo, err := util.GetContextValue[common.SystemInfo](middleware.KEY_JWT_CLAIMS, ctx)
	if err != nil {
		log.Println(err)
	}

	tokenStr, cerr := a.service.SignIn(&signInInput, &systemInfo)
	if cerr.IsError() {
		util.ErrorResponse(ctx, cerr.HttpCode, err.Error())
		return
	}

	token := model.AccessToken{
		Type:  model.Bearer,
		Token: tokenStr,
	}

	util.SuccessResponse(ctx, http.StatusOK, token)
}

func (a authController) Register(ctx *gin.Context) {
	var auth model.SignInInput
	if err := ctx.BindJSON(&auth); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := a.service.SignUp(&auth)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, nil)
}

func (a authController) Logout(ctx *gin.Context) {
	accessClaims, err := util.GetContextValue[model.AccessTokenClaims](middleware.KEY_JWT_CLAIMS, ctx)
	if err != nil {
		log.Println(err)
	}

	if err := a.service.Logout(accessClaims.UserId, accessClaims.RefreshId); err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, nil)
}

func (a authController) LogoutAllDevice(ctx *gin.Context) {
	accessClaims, err := util.GetContextValue[model.AccessTokenClaims](middleware.KEY_JWT_CLAIMS, ctx)
	if err != nil {
		log.Println(err)
	}

	if err := a.service.LogoutAllDevice(accessClaims.UserId); err.IsError() {
		util.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	util.SuccessResponse(ctx, http.StatusOK, nil)
}

func (a authController) RefreshToken(ctx *gin.Context) {
	var token model.AccessToken
	if err := ctx.ShouldBindJSON(&token); err != nil {
		log.Println("Error Refresh Token: ", err)
		util.ErrorResponse(ctx, http.StatusBadRequest, constant.ERR_BODY_REQUEST_MISSING)
		return
	}

	accessToken, err := a.service.RefreshToken(token.Token)
	if err.IsError() {
		util.ErrorResponse(ctx, err.HttpCode, err.Error())
		return
	}

	token = model.AccessToken{
		Type:  "Bearer",
		Token: accessToken,
	}

	util.SuccessResponse(ctx, http.StatusCreated, token)
}
