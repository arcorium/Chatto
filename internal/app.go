package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatto/internal/constant"
	pg_repo "chatto/internal/repository/pg"
	"chatto/internal/repository/redis"
	"chatto/internal/rest/middleware"
	"chatto/internal/ws/manager"

	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/rest"
	"chatto/internal/service"
	"chatto/internal/ws"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewApp(config_ *config.AppConfig) Application {
	return Application{
		Config:        config_,
		App:           gin.New(),
		clientManager: manager.NewClientManager(),
		roomManager:   manager.NewRoomManager(),
	}
}

type Application struct {
	Config *config.AppConfig
	App    *gin.Engine

	clientManager manager.ClientManager
	roomManager   manager.RoomManager
}

func (a *Application) openDatabase() (*gorm.DB, error) {
	conn, err := sql.Open("pgx", a.Config.UserDatabaseURI)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.Credential{}, &model.Room{}, &model.UserRoom{})
	return db, err
}

func (a *Application) openRedisDatabase() (*redis.Client, error) {
	opt := redis.Options{
		Addr:     a.Config.ChatDatabaseURI,
		Username: a.Config.ChatDatabaseUsername,
		Password: a.Config.ChatDatabasePassword,
	}
	client := redis.NewClient(&opt)
	if client == nil {
		return nil, errors.New("redis client is nil")
	}

	status := client.Ping(context.Background())
	return client, status.Err()
}

func (a *Application) setupRedisIndexes(client *redis.Client) error {
	res := client.Ping(context.Background())
	if err := res.Err(); err != nil {
		return err
	}

	//client.Do(context.Background(), "FT.DROPINDEX", constant.REDIS_KEY_USER_INDEX)
	//client.Do(context.Background(), "FT.DROPINDEX", constant.REDIS_KEY_CHAT_INDEX)
	//client.Do(context.Background(), "FT.DROPINDEX", constant.REDIS_KEY_NOTIF_INDEX)

	resInfo := client.Do(context.Background(), "FT.INFO", constant.REDIS_KEY_USER_INDEX)
	if resInfo.Err() != nil {
		// Setup index
		resInfo = client.Do(context.Background(), "FT.CREATE", constant.REDIS_KEY_USER_INDEX, "ON", "HASH", "PREFIX", 1, constant.REDIS_KEY_USER,
			"SCHEMA", "name", "TEXT", "SORTABLE", "role", "TEXT", "SORTABLE", "online", "NUMERIC", "SORTABLE")
		if resInfo.Err() != nil {
			return resInfo.Err()
		}
	}

	resInfo = client.Do(context.Background(), "FT.INFO", constant.REDIS_KEY_CHAT_INDEX)
	if resInfo.Err() != nil {
		// Setup index
		resInfo = client.Do(context.Background(), "FT.CREATE", constant.REDIS_KEY_CHAT_INDEX, "ON", "JSON", "PREFIX", 1, constant.REDIS_KEY_CHAT,
			"SCHEMA", "$.id", "as", "id", "TAG", "$.sender_id", "as", "sender", "TAG", "$.receiver_id", "as", "receiver", "TAG", "$.message", "as", "message", "TEXT", "$.ts", "as", "timestamp", "NUMERIC", "SORTABLE")
		if resInfo.Err() != nil {
			return resInfo.Err()
		}
	}

	resInfo = client.Do(context.Background(), "FT.INFO", constant.REDIS_KEY_NOTIF_INDEX)
	if resInfo.Err() != nil {
		// Setup index
		resInfo = client.Do(context.Background(), "FT.CREATE", constant.REDIS_KEY_NOTIF_INDEX, "ON", "JSON", "PREFIX", 1, constant.REDIS_KEY_NOTIF,
			"SCHEMA", "$.id", "as", "id", "TAG", "$.type", "as", "type", "NUMERIC", "SORTABLE", "$.sender_id", "as", "sender", "TAG", "$.receiver_id", "as", "receiver", "TAG", "$.ts", "as", "timestamp", "NUMERIC", "SORTABLE")
		if resInfo.Err() != nil {
			return resInfo.Err()
		}
	}

	return nil
}

func (a *Application) stopRedis(client *redis.Client) {
	if err := client.Close(); err != nil {
		log.Println(err)
	}
}
func (a *Application) Start() {
	db, err := a.openDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	redisDb, err := a.openRedisDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer a.stopRedis(redisDb)

	err = a.setupRedisIndexes(redisDb)
	if err != nil {
		log.Fatalln(err)
	}

	mw := middleware.NewMiddleware(a.Config)

	userRoomRepo := pg_repo.NewUserRoomRepository(db)

	userRepo := pg_repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	authRepo := pg_repo.NewAuthRepository(db)
	authService := service.NewAuthService(a.Config, authRepo, userService)

	roomRepo := pg_repo.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, userRoomRepo)

	chatRepository := redis_repo.NewChatRepository(redisDb)
	chatService := service.NewChatService(chatRepository, userService, roomService, &a.roomManager, &a.clientManager)

	// Rest Server
	restServer := rest.Server{
		Config:      a.Config,
		Router:      a.App,
		UserService: userService,
		AuthService: authService,
		RoomService: roomService,
		Middleware:  &mw,
	}
	restServer.Setup()

	// Handle Websocket
	wsConfig := ws.WebsocketServerConfig{
		Config:        a.Config,
		Router:        a.App,
		ClientManager: &a.clientManager,
		RoomManager:   &a.roomManager,
		UserService:   userService,
		ChatService:   chatService,
		RoomService:   roomService,
		Middlewares:   &mw,
	}

	wsServer := ws.NewWebsocketServer(&wsConfig)
	wsServer.Setup()

	go func() {
		if err := a.App.Run(a.Config.Address); err != nil {
			log.Fatalln(err)
		}
	}()

	// Create graceful stop
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitChan

	wsServer.Stop()
}
