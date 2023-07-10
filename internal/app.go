package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/repository"
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
		Config: config_,
		App:    gin.New(),
	}
}

type Application struct {
	Config *config.AppConfig
	App    *gin.Engine
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

	err = db.AutoMigrate(&model.User{}, &model.TokenDetails{}, &model.Room{})
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

func (a *Application) Start() {
	db, err := a.openDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	redisDb, err := a.openRedisDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer func(redisDb *redis.Client) {
		err := redisDb.Close()
		if err != nil {
			log.Println(err)
		}
	}(redisDb)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(a.Config, authRepo, userRepo)

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)

	chatRepository := repository.NewChatRepository(redisDb)
	chatService := service.NewChatService(chatRepository)

	// Rest Server
	restServer := rest.Server{
		Config:      a.Config,
		Router:      a.App,
		UserService: userService,
		AuthService: authService,
		RoomService: roomService,
	}
	restServer.Setup()

	// Handle Websocket routes
	wsServer := ws.NewWebsocketServer(a.Config, a.App, userService, authService, chatService)
	wsServer.Setup()

	log.Fatalln(a.App.Run(a.Config.Address))
}
