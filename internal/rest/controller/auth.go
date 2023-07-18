package controller

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/util/httputil"
	"chatto/internal/util/strutil"

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
	authRoute := router.Group("/auth", middlewares.UserAgent.Handle())
	authRoute.POST("/login", a.Login)
	authRoute.POST("/register", a.Register)
	authRoute.POST("/refresh", a.RefreshToken)

	authRoute.Use(middlewares.TokenValidation.Handle())
	authRoute.GET("/devices/", a.GetLoginDevice)
	authRoute.POST("/logout", a.Logout)
	authRoute.POST("/logout/:id", a.LogoutDevice)
	authRoute.POST("/logout/all", a.LogoutAllDevice)
}

func (a authController) Login(ctx *gin.Context) {
	var signInInput dto.SignInInput
	if err := ctx.BindJSON(&signInInput); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	systemInfo, err := util.GetContextValue[common.SystemInfo](constant.KEY_USER_AGENT, ctx)
	if err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.USER_AGENT_UNKNOWN_ERROR, constant.MSG_UNKNOWN_USER_AGENT))
		return
	}

	output, cerr := a.service.SignIn(&signInInput, systemInfo)
	httputil.ConditionalResponse(ctx, cerr, http.StatusInternalServerError, http.StatusOK, output)
}

func (a authController) Register(ctx *gin.Context) {
	var auth dto.SignUpInput
	if err := ctx.BindJSON(&auth); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}

	err := a.service.SignUp(&auth)
	httputil.ConditionalResponse(ctx, err, http.StatusInternalServerError, http.StatusOK, nil)
}

func (a authController) Logout(ctx *gin.Context) {
	accessClaims, err := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)
	if err != nil {
		log.Println(err)
	}

	cerr := a.service.Logout(accessClaims.UserId, accessClaims.RefreshId)
	httputil.ConditionalResponse(ctx, cerr, http.StatusInternalServerError, http.StatusOK, nil)
}

func (a authController) LogoutAllDevice(ctx *gin.Context) {
	accessClaims, err := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)
	if err != nil {
		log.Println(err)
	}

	cerr := a.service.LogoutAllDevice(accessClaims.UserId)
	httputil.ConditionalResponse(ctx, cerr, http.StatusUnauthorized, http.StatusOK, nil)
}

func (a authController) RefreshToken(ctx *gin.Context) {
	var token dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&token); err != nil {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_BODY_REQUEST_ERROR, constant.MSG_BAD_BODY_REQUEST))
		return
	}
	output, err := a.service.RefreshToken(&token)
	httputil.ConditionalResponse(ctx, err, http.StatusUnauthorized, http.StatusCreated, output)
}

func (a authController) GetLoginDevice(ctx *gin.Context) {
	claims, _ := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)

	sysInfos, err := a.service.GetLoginDevices(claims.UserId)

	httputil.ConditionalResponse(ctx, err, http.StatusNotFound, http.StatusOK, sysInfos)
}

func (a authController) LogoutDevice(ctx *gin.Context) {
	refreshId := ctx.Param("id")
	if strutil.IsEmpty(refreshId) {
		httputil.ErrorResponse(ctx, http.StatusBadRequest, common.NewError(common.BAD_PARAMETER_ERROR, constant.MSG_URI_PARAM_MISSING))
		return
	}

	claims, _ := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)

	cerr := a.service.Logout(claims.UserId, refreshId)
	httputil.ConditionalResponse(ctx, cerr, http.StatusInternalServerError, http.StatusOK, nil)
}
