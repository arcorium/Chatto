package ws

import (
	"fmt"

	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/ws/controller"
	"chatto/internal/ws/handler"
	"chatto/internal/ws/manager"

	"github.com/gin-gonic/gin"
)

type WebsocketServerConfig struct {
	Config *config.AppConfig
	Router gin.IRouter

	ClientManager *manager.ClientManager
	RoomManager   *manager.RoomManager

	UserService service.IUserService
	ChatService service.IChatService
	RoomService service.IRoomService

	Middlewares *middleware.Middleware
}

func NewWebsocketServer(config *WebsocketServerConfig) Server {
	return Server{
		cfg:           config.Config,
		router:        config.Router,
		clientManager: config.ClientManager,
		roomManager:   config.RoomManager,
		payloadChan:   make(chan *model.Payload, 100),
		clientChan:    make(chan *model.Client),
		userService:   config.UserService,
		chatService:   config.ChatService,
		roomService:   config.RoomService,
		middlewares:   config.Middlewares,
	}
}

type Server struct {
	cfg    *config.AppConfig
	router gin.IRouter

	clientManager *manager.ClientManager
	roomManager   *manager.RoomManager
	payloadChan   chan *model.Payload
	clientChan    chan *model.Client

	userService service.IUserService
	chatService service.IChatService
	roomService service.IRoomService

	middlewares *middleware.Middleware
}

// lookupRooms Should be called when chat service start, it will get all the rooms from room service and create appropriate ChatRoom
func (s *Server) lookupRooms() error {
	rooms, err := s.roomService.FindRooms()
	if err.IsError() {
		return err.Error()
	}

	for _, room := range rooms {
		var chatRoom model.ChatRoom
		if room.Private {
			chatRoom = model.NewPrivateChatRoom(room.Id, room.Name, room.Description, room.InviteOnly)
		} else {
			chatRoom = model.NewChatRoom(room.Id, room.Name, room.Description, room.InviteOnly)
		}
		s.roomManager.AddRooms(&chatRoom)
	}
	return nil
}

func (s *Server) Setup() {
	handler.StartClientHandler(s.chatService, s.clientChan, s.payloadChan)
	handler.StartPayloadHandler(s.payloadChan, s.chatService, s.roomManager, s.clientManager)

	if err := s.lookupRooms(); err != nil {
		panic(fmt.Sprint("Error on lookupRooms: ", err))
	}
	// Set redis indexes

	websocketHandler := controller.NewWebsocketHandler(s.clientChan)
	websocketHandler.Route(s.router, s.middlewares)
}

func (s *Server) Stop() {
	close(s.payloadChan)
	close(s.clientChan)

	s.chatService.ClearUsers()
	s.clientManager.StopClientChannels()
}
