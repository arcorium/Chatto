package route

import (
	"chatto/internal/config"
	"chatto/internal/rest/controller"
	"chatto/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

func NewAuthRoute(authController controller.IAuthController) AuthRoute {
	return AuthRoute{authController: authController}
}

type AuthRoute struct {
	authController controller.IAuthController
}

func (a AuthRoute) V1Handle(router gin.IRouter, cfg *config.AppConfig) {
	userAgentMiddleware := middleware.UserAgentValidationMiddleware{}
	authMiddleware := middleware.TokenValidationMiddleware{Config: &middleware.TokenValidationConfig{
		SecretKeyFunc: cfg.JWTKeyFunc,
		TokenType:     "Bearer",
		SigningType:   "HS512",
	}}

	authRoute := router.Group("/auth")

	authRoute.POST("/login", userAgentMiddleware.Handle(), a.authController.Login)
	authRoute.POST("/register", userAgentMiddleware.Handle(), a.authController.Register)

	authRoute.POST("/refresh", a.authController.RefreshToken)

	authRoute.Use(authMiddleware.Handle())
	authRoute.POST("/logout", a.authController.Logout)
	authRoute.POST("/logout/all", a.authController.LogoutAllDevice)
}
