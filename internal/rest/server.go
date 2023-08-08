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
	Middleware  *middleware.Middleware
}

func (s *Server) registerControllers(controllers ...controller.IController) {
	v1Router := s.Router.Group("/api/v1")
	for _, c := range controllers {
		c.Route(v1Router, s.Middleware)
	}
}

func (s *Server) Setup() {
	userController := controller.NewUserController(s.UserService)
	authController := controller.NewAuthController(s.AuthService)
	roomController := controller.NewRoomController(s.RoomService)

	// Handle REST API routes
	s.registerControllers(userController, authController, roomController)
}
