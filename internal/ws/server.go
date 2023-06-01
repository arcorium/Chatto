package ws

import (
	"server_client_chat/internal/config"
	"server_client_chat/internal/model"
	"server_client_chat/internal/service"
	"server_client_chat/internal/ws/handler"
	"server_client_chat/internal/ws/manager"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

func NewWebsocketServer(config *config.AppConfig, router gin.IRouter, userService service.IUserService, authService service.IAuthService, chatService service.IChatService) Server {
	return Server{
		cfg:           config,
		router:        router,
		clientManager: manager.NewClientManager(),
		roomManager:   manager.NewRoomManager(),
		payloadChan:   make(chan *model.Payload, 100),
		clientChan:    make(chan *model.Client),
		userService:   userService,
		authService:   authService,
		chatService:   chatService,
	}
}

type Server struct {
	cfg    *config.AppConfig
	router gin.IRouter

	clientManager manager.ClientManager
	roomManager   manager.RoomManager
	payloadChan   chan *model.Payload
	clientChan    chan *model.Client

	userService service.IUserService
	authService service.IAuthService
	chatService service.IChatService
}

func (s *Server) Init() {
	handler.StartClientHandler(s.chatService, &s.clientManager, s.clientChan, s.payloadChan)
	handler.StartPayloadHandler(&s.clientManager, &s.roomManager, s.payloadChan)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	websocketHandler := handler.NewWebsocketHandler(upgrader, s.clientChan)
	route := NewWebsocketRoute(websocketHandler)
	route.RegisterRoute(s.cfg, s.router)
}
