package rest

import (
	"server_client_chat/internal/config"
	"server_client_chat/internal/rest/controller"
	"server_client_chat/internal/rest/route"
	"server_client_chat/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *config.AppConfig
	Router gin.IRouter

	UserService service.IUserService
	AuthService service.IAuthService
}

func (s *Server) registerRoutes(routes ...route.IRoute) {
	route.V1Route(s.Router, s.Config, routes...)
}

func (s *Server) Setup() {

	userController := controller.NewUserController(s.UserService)
	userRoute := route.NewUserRoute(userController)

	authController := controller.NewAuthController(s.AuthService)
	authRoute := route.NewAuthRoute(&authController)

	// Handle REST API routes
	s.registerRoutes(authRoute, userRoute)
}
