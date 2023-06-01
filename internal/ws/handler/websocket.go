package handler

import (
	"log"

	"server_client_chat/internal/model"
	"server_client_chat/internal/rest/middleware"
	"server_client_chat/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewWebsocketHandler(upgrader websocket.Upgrader, client chan<- *model.Client) WebsocketHandler {
	return WebsocketHandler{upgrader: upgrader, client: client}
}

type WebsocketHandler struct {
	upgrader websocket.Upgrader
	client   chan<- *model.Client
}

func (w *WebsocketHandler) ServeWebsocket(ctx *gin.Context) {

	conn, err := w.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
	}

	// Handle new client
	details, _ := util.GetContextValue[model.AccessTokenClaims](middleware.KEY_JWT_CLAIMS, ctx)

	client := model.NewClient(details.UserId, model.ClientStatusRegister, conn)
	w.RegisterNewClient(&client)
	// Let the connection done
}

func (w *WebsocketHandler) RegisterNewClient(client *model.Client) {
	w.client <- client
}
