package handler

import (
	"log"
	"net/http"

	"chatto/internal/model"
	"chatto/internal/rest/middleware"
	"chatto/internal/service"
	"chatto/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewWebsocketHandler(upgrader websocket.Upgrader, userService service.IUserService, client chan<- *model.Client) WebsocketHandler {
	return WebsocketHandler{upgrader: upgrader, userService: userService, client: client}
}

type WebsocketHandler struct {
	upgrader websocket.Upgrader

	userService service.IUserService

	client chan<- *model.Client
}

func (w *WebsocketHandler) ServeWebsocket(ctx *gin.Context) {

	conn, err := w.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
	}

	// Handle new client
	details, _ := util.GetContextValue[model.AccessTokenClaims](middleware.KEY_JWT_CLAIMS, ctx)

	user, cerr := w.userService.FindUserById(details.UserId)
	if cerr.IsError() {
		model.NewErrorResponse(http.StatusUnauthorized, err.Error(), nil)
		return
	}

	client := model.NewClient(user.UserId, user.Username, user.Role, model.ClientStatusRegister, conn)
	w.RegisterNewClient(&client)
	// Let the connection done
}

func (w *WebsocketHandler) RegisterNewClient(client *model.Client) {
	w.client <- client
}
