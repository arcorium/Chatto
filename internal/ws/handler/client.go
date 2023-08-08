package handler

import (
	"log"

	"chatto/internal/model"
	"chatto/internal/service"

	"github.com/gorilla/websocket"
)

func StartClientHandler(chatService service.IChatService, client <-chan *model.Client, payload chan<- *model.Payload) {
	handler := &ClientHandler{
		chatService: chatService,
		payload:     payload,
		client:      client,
	}

	go handler.ClientHandle()
}

type ClientHandler struct {
	payload chan<- *model.Payload
	client  <-chan *model.Client

	chatService service.IChatService
}

func (c *ClientHandler) ClientHandle() {
	for client := range c.client {
		if c.registerClient(client) != nil {
			c.unregisterClient(client)
			continue
		}
		go c.ClientReadHandle(client)
		go c.ClientWriteHandle(client)
	}

	log.Println("Client Stopped")
}

// ClientReadHandle Handle for reading each client message
func (c *ClientHandler) ClientReadHandle(client *model.Client) {
	defer func() {
		if err := client.Conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		var input model.PayloadInput
		err := client.Conn.ReadJSON(&input)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, 10054) {
				c.unregisterClient(client)
				break
			}
		}
		payload := c.chatService.ProcessPayload(client, &input)
		c.payload <- &payload
	}
}

// ClientWriteHandle Handle for writing message to each client
func (c *ClientHandler) ClientWriteHandle(client *model.Client) {
	for msg := range client.IncomingPayload {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				break
			}
			log.Println("Client ", *client, " Write Error: ", err)
			continue
		}
	}
}

func (c *ClientHandler) registerClient(client *model.Client) error {
	log.Println("Registering client: ", client)
	cerr := c.chatService.NewClient(client)
	return cerr.Error()
}

func (c *ClientHandler) unregisterClient(client *model.Client) {
	log.Println("Unregistering client: ", client)
	_ = c.chatService.RemoveClient(client)
}
