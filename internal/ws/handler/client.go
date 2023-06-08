package handler

import (
	"chatto/internal/constant"
	"log"
	"time"

	"chatto/internal/model"
	"chatto/internal/service"
	"chatto/internal/ws/manager"

	"github.com/gorilla/websocket"
)

func StartClientHandler(chatService service.IChatService, clientManager *manager.ClientManager, client <-chan *model.Client, payload chan<- *model.Payload) {
	handler := &ClientHandler{
		service:       chatService,
		clientManager: clientManager,
		payload:       payload,
		client:        client,
	}

	go handler.ClientHandle()
}

type ClientHandler struct {
	service       service.IChatService
	clientManager *manager.ClientManager
	payload       chan<- *model.Payload
	client        <-chan *model.Client
}

func (c *ClientHandler) ClientHandle() {
	for {
		select {
		case client := <-c.client:
			switch client.Status {
			case model.ClientStatusRegister:
				c.registerClient(client)
				go c.ClientReadHandle(client)
				go c.ClientWriteHandle(client)
				//case model.ClientStatusUnregister:
				//	c.unregisterClient(client)
			}
		}
	}
}

// ClientReadHandle Handle for reading each client message
func (c *ClientHandler) ClientReadHandle(client *model.Client) {
	defer func() {
		if err := client.Conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	// TODO: Handle read timeout
	client.Conn.SetReadLimit(constant.CLIENT_READ_LIMIT_SIZE)
	err := client.Conn.SetReadDeadline(time.Now().Add(constant.CLIENT_READ_LIMIT_TIME))
	if err != nil {
		log.Println(err)
		return
	}

	for {
		var payload model.Payload
		err = client.Conn.ReadJSON(&payload)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, 10054) {
				client.Status = model.ClientStatusUnregister
				c.unregisterClient(client)
				break
			}
			// Response error to client
			payload = model.NewErrorPayload(constant.ERR_BAD_PAYLOAD)
		}
		payload.Populate(client)
		c.payload <- &payload
	}
}

// ClientWriteHandle Handle for writing message to each client
func (c *ClientHandler) ClientWriteHandle(client *model.Client) {
	// TODO: Use ticker to detect if there is respond from client, otherwise return
	for {
		select {
		case msg := <-client.IncomingPayload:
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
}

func (c *ClientHandler) registerClient(client *model.Client) {
	client.Status = model.ClientStatusOnline
	c.clientManager.AddClients(client)
	log.Println("Registering client: ", client)
	// TODO: Add service handle
}

func (c *ClientHandler) unregisterClient(client *model.Client) {
	client.Status = model.ClientStatusOffline
	c.clientManager.RemoveClientById(client.Id)
	log.Println("Unregistering client: ", client)
	// TODO: Add service handle
}
