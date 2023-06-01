package handler

import (
	"log"

	"server_client_chat/internal/model"
	"server_client_chat/internal/ws/manager"
)

func StartPayloadHandler(clientManager *manager.ClientManager, roomManager *manager.RoomManager, payload <-chan *model.Payload) {
	handler := &PayloadHandler{
		clientManager: clientManager,
		roomManager:   roomManager,
		payload:       payload,
	}

	go handler.PayloadHandle()
}

type PayloadHandler struct {
	roomManager   *manager.RoomManager
	clientManager *manager.ClientManager
	payload       <-chan *model.Payload
}

func (p *PayloadHandler) PayloadHandle() {
	for {
		select {
		case payload := <-p.payload:
			log.Println(payload)
			// Convert interface{} into bytes, so it can be unmarshalled
			bytes, err := payload.EncodeData()
			if err != nil {
				log.Println(err)
				continue
			}
			// Get client
			client, err := p.clientManager.GetClientById(payload.ClientId)
			if err != nil {
				log.Println(err)
				continue
			}

			// Check payload type
			switch payload.Type {
			case model.PayloadMessage:
				msg, err := model.Decode[model.Message](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				msg.Populate()
				p.HandleMessage(client, &msg)
			case model.PayloadNotification:
				notif, err := model.Decode[model.Notification](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				notif.Populate()
				p.HandleNotification(client, &notif)
			case model.PayloadStartChat:
				privateChat, err := model.Decode[model.PrivateChat](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandlePrivateRoom(client, &privateChat)
				// Let client handle it
			case model.PayloadCreateRoom:
				createRoom, err := model.Decode[model.CreateRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleCreateRoom(client, &createRoom)
			case model.PayloadJoinRoom:
				joinRoom, err := model.Decode[model.JoinRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleJoinRoom(client, &joinRoom)
			case model.PayloadLeaveRoom:
				// Implicitly remove room when there is only one user there
				leaveRoom, err := model.Decode[model.JoinRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleLeaveRoom(client, &leaveRoom)
			}
		}
	}
}

// HandleMessage TODO: Move into another handler
func (p *PayloadHandler) HandleMessage(client *model.Client, message *model.Message) {
	// Find Room by the room_id
	room, err := p.roomManager.GetRoomById(message.Receiver)
	if err != nil {
		log.Println(err)
		return
	}
	// Send message into all clients in room
	payload := model.NewMessagePayload(client, message)
	room.BroadcastPayloadExceptClientId(payload, client.Id)
}

// HandleNotification TODO: Move into another handler
func (p *PayloadHandler) HandleNotification(client *model.Client, notification *model.Notification) {
	// TODO: Broadcast notification to all clients
	p.sendRoomNotification(client, notification)
}

func (p *PayloadHandler) sendRoomNotification(client *model.Client, notification *model.Notification) {
	// Find Room by the room_id
	room, err := p.roomManager.GetRoomById(notification.Receiver)
	if err != nil {
		log.Println(err)
		return
	}
	// Send message into all clients in room
	payload := model.NewNotificationPayload(notification)
	room.BroadcastPayloadExceptUserId(payload, client.UserId)
}

func (p *PayloadHandler) HandlePrivateRoom(client *model.Client, chat *model.PrivateChat) {
	// Get sender and opponent client and check if the opponent is online, when offline just add it on the redis so the opponent will get the chat when online
	clients := p.clientManager.GetClientsByUserId(chat.Opponent)
	if len(clients) == 0 {

	}
	clients = append(clients, client)

	// Create room
	room := model.NewRoom(chat.Opponent, true, clients...)
	p.roomManager.AddRooms(&room)

	// Respond with room_id
	client.IncomingPayload <- model.NewRoomPayload(&room)
}

func (p *PayloadHandler) HandleCreateRoom(client *model.Client, createRoom *model.CreateRoom) {
	room := model.NewRoom(createRoom.Name, createRoom.Private, client)

	// Get each client for members, doing the same in private chat when the member is offline
	// Client
	room.AddClients(p.clientManager.GetClientsByUserId(client.UserId)...)
	// Members
	for _, m := range createRoom.Members {
		clients := p.clientManager.GetClientsByUserId(m)
		room.AddClients(clients...)
	}
	p.roomManager.AddRooms(&room)

	// respond with room_id
	client.IncomingPayload <- model.NewRoomPayload(&room)
}

func (p *PayloadHandler) HandleJoinRoom(client *model.Client, joinRoom *model.JoinRoom) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		payload := model.NewErrorResponsePayload(err.Error())
		client.SendPayload(&payload)
		return
	}
	if room.Private {
		payload := model.NewErrorResponsePayload("cannot join private room")
		client.SendPayload(&payload)
		return
	}
	room.AddClients(p.clientManager.GetClientsByUserId(client.UserId)...)
	// TODO: Using username instead of user id
	payload := model.NewNotificationPayload(model.NewServerNotification(room.Id, model.NotifJoinRoom, client.UserId+" join the chat"))
	room.BroadcastPayload(payload)
}

func (p *PayloadHandler) HandleLeaveRoom(client *model.Client, joinRoom *model.JoinRoom) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		payload := model.NewErrorResponsePayload(err.Error())
		client.SendPayload(&payload)
		return
	}
	room.RemoveClientsByUserId(client.UserId)

	// TODO: Using username instead of user id
	payload := model.NewNotificationPayload(model.NewServerNotification(room.Id, model.NotifJoinRoom, client.UserId+" leave the chat"))
	room.BroadcastPayload(payload)
}
