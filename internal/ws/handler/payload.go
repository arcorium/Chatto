package handler

import (
	"chatto/internal/constant"
	"log"

	"chatto/internal/model"
	"chatto/internal/ws/manager"
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
			case model.PayloadNotification:
				notif, err := model.Decode[model.Notification](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				notif.Populate()
				p.HandleNotification(client, &notif)
			case model.PayloadPrivateChat:
				privateChat, err := model.Decode[model.IncomeMessage](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandlePrivateChat(client, &privateChat)
				// Let client handle it
			case model.PayloadRoomChat:
				roomChat, err := model.Decode[model.IncomeMessage](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				if roomChat.Receiver == "general" {
					p.HandleGeneralChat(client, &roomChat)
				} else {
					p.HandleRoomChat(client, &roomChat)
				}
			case model.PayloadCreateRoom:
				createRoom, err := model.Decode[model.CreateRoomPayload](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleCreateRoom(client, &createRoom)
			case model.PayloadJoinRoom:
				joinRoom, err := model.Decode[model.JoinRoomPayload](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleJoinRoom(client, &joinRoom)
			case model.PayloadLeaveRoom:
				// Implicitly remove room when there is only one user there
				leaveRoom, err := model.Decode[model.JoinRoomPayload](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleLeaveRoom(client, &leaveRoom)
			case model.PayloadGetUsers:
				getUser, err := model.Decode[model.GetUserPayload](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleGetUsers(client, &getUser)
			}
		}
	}
}

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
	payload := model.NewNotificationPayload(client, notification)
	room.BroadcastPayloadExceptUserId(payload, client.Id)
}

func (p *PayloadHandler) HandlePrivateChat(client *model.Client, chat *model.IncomeMessage) {
	// Get sender and opponent client and check if the opponent is online, when offline just add it on the redis so the opponent will get the chat when online
	receivers := p.clientManager.GetClientsByUserId(chat.Receiver)
	if len(receivers) == 0 {
		payload := model.NewErrorPayload(constant.ERR_CLIENT_NOT_EXIST)
		client.SendPayload(&payload)
		return
	}
	senders := p.clientManager.GetClientsByUserId(client.UserId)
	receivers = append(receivers, senders...)

	// Respond with room_id
	message := model.NewMessage(chat)
	payload := model.NewPrivateMessagePayload(client, message)
	for _, c := range receivers {
		if c.Id == client.Id {
			continue
		}
		c.SendPayload(&payload)
	}
}

func (p *PayloadHandler) HandleRoomChat(client *model.Client, chat *model.IncomeMessage) {
	// Search rooms
	room, err := p.roomManager.GetRoomById(chat.Receiver)
	if err != nil {
		payload := model.NewErrorPayload(constant.ERR_ROOM_NOT_EXIST)
		client.SendPayload(&payload)
		return
	}

	message := model.NewMessage(chat)

	// Broadcast
	payload := model.NewRoomMessagePayload(client, message)
	room.BroadcastPayloadExceptClientId(&payload, client.Id)
}

func (p *PayloadHandler) HandleGeneralChat(client *model.Client, chat *model.IncomeMessage) {
	message := model.NewMessage(chat)
	payload := model.NewRoomMessagePayload(client, message)
	for _, c := range p.clientManager.Clients {
		if c.Id == client.Id {
			continue
		}
		c.SendPayload(&payload)
	}
}

func (p *PayloadHandler) HandleCreateRoom(client *model.Client, createRoom *model.CreateRoomPayload) {
	room := model.NewRoom(createRoom.Name, createRoom.Private, client)

	// Get each client for members, doing the same in private chat when the member is offline
	// Client
	room.AddClients(p.clientManager.GetClientsByUserId(client.Id)...)
	// MemberIds
	for _, m := range createRoom.MemberIds {
		clients := p.clientManager.GetClientsByUserId(m)
		room.AddClients(clients...)
	}
	p.roomManager.AddRooms(&room)

	// respond with room_id
	client.SendPayload(model.NewRoomPayload(&room))
}

func (p *PayloadHandler) HandleJoinRoom(client *model.Client, joinRoom *model.JoinRoomPayload) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		payload := model.NewErrorPayload(err.Error())
		client.SendPayload(&payload)
		return
	}
	if room.Private {
		payload := model.NewErrorPayload("cannot join private room")
		client.SendPayload(&payload)
		return
	}
	room.AddClients(p.clientManager.GetClientsByUserId(client.UserId)...)
	// TODO: Using username instead of user id
	notif := model.NewNotification(room.Id, model.NotifJoinRoom)
	payload := model.NewNotificationPayload(client, notif)
	room.BroadcastPayload(payload)
}

func (p *PayloadHandler) HandleLeaveRoom(client *model.Client, joinRoom *model.JoinRoomPayload) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		payload := model.NewErrorPayload(err.Error())
		client.SendPayload(&payload)
		return
	}
	room.RemoveClientsByUserId(client.UserId)

	notif := model.NewNotification(room.Id, model.NotifLeaveRoom)
	payload := model.NewNotificationPayload(client, notif)
	room.BroadcastPayload(payload)
}

func (p *PayloadHandler) HandleGetUsers(client *model.Client, userPayload *model.GetUserPayload) {
	clients := p.clientManager.GetClientsByUsername(userPayload.Username)
	// Send back with the users
	payload := model.NewUserResponsePayload(client, clients)
	client.SendPayload(&payload)
}
