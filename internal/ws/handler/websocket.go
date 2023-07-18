package handler

import (
	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/service"
	"chatto/internal/util"
	"chatto/internal/util/httputil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewWebsocketHandler(upgrader websocket.Upgrader, userService service.IUserService, client chan<- *model.Client) WebsocketHandler {
	return WebsocketHandler{upgrader: upgrader, userService: userService, client: client}
}

type WebsocketHandler struct {
	upgrader websocket.Upgrader

	userService service.IUserService
	client      chan<- *model.Client
}

func (w *WebsocketHandler) ServeWebsocket(ctx *gin.Context) {

	conn, err := w.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
	}

	details, _ := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)

	user, cerr := w.userService.FindUserById(details.UserId)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, http.StatusUnauthorized, cerr)
	} else {
		client := model.NewClient(user.Id, user.Username, user.Role, conn)
		w.registerClient(&client)
	}
}

func (w *WebsocketHandler) registerClient(client *model.Client) {
	w.client <- client
}
