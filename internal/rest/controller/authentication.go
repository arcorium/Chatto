package controller

import (
	"chatto/internal/constant"
	"log"
	"net/http"

	"chatto/internal/model"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/util"

	"github.com/gin-gonic/gin"
)

func NewAuthController(service service.IAuthService) AuthController {
	return AuthController{service: service}
}

type AuthController struct {
	service service.IAuthService
}

func (a *AuthController) Login(ctx *gin.Context) {
	var signInInput model.SignInInput
	if err := ctx.BindJSON(&signInInput); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	systemInfo, err := util.GetContextValue[model.SystemInfo](middleware.KEY_JWT_CLAIMS, ctx)
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

func (a *AuthController) Register(ctx *gin.Context) {
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

func (a *AuthController) Logout(ctx *gin.Context) {
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

func (a *AuthController) LogoutAllDevice(ctx *gin.Context) {
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

func (a *AuthController) RefreshToken(ctx *gin.Context) {
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
