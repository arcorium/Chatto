package route

import (
	"chatto/internal/config"
	"chatto/internal/rest/controller"
	"chatto/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

func NewUserRoute(userController controller.IUserController) UserRoute {
	return UserRoute{userController: userController}
}

type UserRoute struct {
	userController controller.IUserController
}

func (u UserRoute) V1Handle(router gin.IRouter, cfg *config.AppConfig) {
	authMiddleware := middleware.TokenValidationMiddleware{Config: &middleware.TokenValidationConfig{
		SecretKeyFunc: cfg.JWTKeyFunc,
		TokenType:     "Bearer",
		SigningType:   "HS512",
	}}

	userRoute := router.Group("/users", authMiddleware.Handle())
	userRoute.GET("/", u.userController.GetUsers)
	userRoute.GET("/:id", u.userController.GetUserById)
	userRoute.PUT("/:id", u.userController.UpdateUser)
	userRoute.DELETE("/:id", u.userController.RemoveUser)

}
