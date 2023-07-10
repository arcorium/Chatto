package rest

import (
	"chatto/internal/config"
	"chatto/internal/rest/controller"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *config.AppConfig
	Router gin.IRouter

	UserService service.IUserService
	AuthService service.IAuthService
	RoomService service.IRoomService
}

func (s *Server) registerControllers(controllers ...controller.IController) {
	middlewares := middleware.NewMiddleware(s.Config)
	for _, c := range controllers {
		c.Route(s.Router, &middlewares)
	}
}

func (s *Server) Setup() {
	userController := controller.NewUserController(s.UserService)
	authController := controller.NewAuthController(s.AuthService)
	roomController := controller.NewRoomController(s.RoomService)

	// Handle REST API routes
	s.registerControllers(userController, authController, roomController)
}
