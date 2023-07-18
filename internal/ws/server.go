package ws

import (
	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/service"
	"chatto/internal/ws/handler"
	"chatto/internal/ws/manager"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

func NewWebsocketServer(config *config.AppConfig, router gin.IRouter, userService service.IUserService, roomService service.IRoomService, chatService service.IChatService) Server {
	return Server{
		cfg:           config,
		router:        router,
		clientManager: manager.NewClientManager(),
		roomManager:   manager.NewRoomManager(),
		payloadChan:   make(chan *model.Payload, 100),
		clientChan:    make(chan *model.Client),
		userService:   userService,
		roomService:   roomService,
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
	chatService service.IChatService
	roomService service.IRoomService
}

func (s *Server) Setup() {
	handler.StartClientHandler(s.chatService, &s.clientManager, s.clientChan, s.payloadChan)
	handler.StartPayloadHandler(&s.clientManager, &s.roomManager, s.payloadChan, s.chatService, s.roomService)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	websocketHandler := handler.NewWebsocketHandler(upgrader, s.userService, s.clientChan)
	route := NewWebsocketRoute(websocketHandler)
	route.RegisterRoute(s.cfg, s.router)
}
