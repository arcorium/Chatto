package controller

import (
	"log"
	"net/http"

	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/rest/controller"
	"chatto/internal/rest/middleware"
	"chatto/internal/util"
	"chatto/internal/util/httputil"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewWebsocketHandler(client chan<- *model.Client) controller.IController {
	return &WebsocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		client: client,
	}
}

type WebsocketHandler struct {
	upgrader websocket.Upgrader

	client chan<- *model.Client
}

func (w *WebsocketHandler) Route(router gin.IRouter, middleware *middleware.Middleware) {
	router.GET("/chat", middleware.TokenValidation, w.ServeWebsocket)
}

func (w *WebsocketHandler) ServeWebsocket(ctx *gin.Context) {
	claims, err := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, ctx)

	if err != nil {
		httputil.ErrorResponse(ctx, http.StatusUnauthorized, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR))
	} else {
		conn, err := w.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println(err)
		}

		client := model.NewClient(claims.UserId, claims.Name, claims.Role, conn)
		w.registerClient(&client)
	}
}

func (w *WebsocketHandler) registerClient(client *model.Client) {
	w.client <- client
}
