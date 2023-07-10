package handler

import (
	"log"

	"chatto/internal/model"
	"chatto/internal/service"
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
	chatService   service.IChatService
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
			case model.PayloadPrivateNotification:
				notif, err := model.Decode[model.IncomeNotification](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandlePrivateNotification(client, &notif)
			case model.PayloadRoomNotification:
				notif, err := model.Decode[model.IncomeNotification](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleRoomNotification(client, &notif)
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
				p.HandleRoomChat(client, &roomChat)
			case model.PayloadCreateRoom:
				createRoom, err := model.Decode[model.IncomeCreateRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleCreateRoom(client, &createRoom)
			case model.PayloadJoinRoom:
				joinRoom, err := model.Decode[model.IncomeJoinRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleJoinRoom(client, &joinRoom)
			case model.PayloadLeaveRoom:
				// Implicitly remove room when there is only one user there
				leaveRoom, err := model.Decode[model.IncomeJoinRoom](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleLeaveRoom(client, &leaveRoom)
			case model.PayloadGetUsers:
				getUser, err := model.Decode[model.IncomeGetUser](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleGetUsers(client, &getUser)
			}
		}
	}
}

func (p *PayloadHandler) HandlePrivateNotification(sender *model.Client, notification *model.IncomeNotification) {
	receivers := p.clientManager.GetClientsByUserId(notification.Receiver)
	if len(receivers) == 0 {
		return
	}

	privateNotif, err := p.chatService.HandleNewPrivateNotification(sender, notification)
	if err.IsError() {
		return
	}

	senders := p.clientManager.GetUniqueClientByUserId(sender)
	if len(senders) >= 1 {
		forwardNotification := model.NewOutcomeNotificationForward(receivers[0], &privateNotif)
		forwardPayload := model.NewPayload(model.PayloadForwarder, &forwardNotification)
		for _, c := range senders {
			c.SendPayload(&forwardPayload)
		}
	}

	privateNotifPayload := model.NewPayload(model.PayloadPrivateChat, &privateNotif)
	for _, c := range receivers {
		c.SendPayload(&privateNotifPayload)
	}
}

func (p *PayloadHandler) HandleRoomNotification(sender *model.Client, notification *model.IncomeNotification) {
	// Search rooms
	room, err := p.roomManager.GetRoomById(notification.Receiver)
	if err != nil {
		return
	}

	roomNotification, err := p.chatService.HandleNewRoomNotification(sender, notification)
	if err != nil {
		return
	}

	// Broadcast
	payload := model.NewPayload(model.PayloadRoomNotification, &roomNotification)
	room.Broadcast(&payload, sender.Id)
}

func (p *PayloadHandler) HandlePrivateChat(sender *model.Client, chat *model.IncomeMessage) {
	receivers := p.clientManager.GetClientsByUserId(chat.Receiver)
	if len(receivers) == 0 {
		return
	}

	privateMessage, err := p.chatService.HandleNewPrivateMessage(sender, chat)
	if err.IsError() {
		return
	}

	// Send to all clients for the same user id
	senders := p.clientManager.GetUniqueClientByUserId(sender)
	if len(senders) >= 1 {
		forwardMessage := model.NewOutcomeMessageForward(receivers[0], &privateMessage)
		forwardPayload := model.NewPayload(model.PayloadForwarder, &forwardMessage)
		for _, c := range senders {
			c.SendPayload(&forwardPayload)
		}
	}

	privateMessagePayload := model.NewPayload(model.PayloadPrivateChat, &privateMessage)
	for _, c := range receivers {
		c.SendPayload(&privateMessagePayload)
	}
}

func (p *PayloadHandler) HandleRoomChat(sender *model.Client, chat *model.IncomeMessage) {
	// Search rooms
	// TODO: Room is saved on redis, so each application starts it should take from redis
	room, err := p.roomManager.GetRoomById(chat.Receiver)
	if err != nil {
		return
	}

	roomMessage, err := p.chatService.HandleNewRoomMessage(sender, chat)
	if err != nil {
		return
	}

	// Broadcast
	payload := model.NewPayload(model.PayloadRoomChat, &roomMessage)
	room.Broadcast(&payload, sender.Id)
}

func (p *PayloadHandler) HandleCreateRoom(client *model.Client, createRoom *model.IncomeCreateRoom) {
	outcome, err := p.chatService.HandleNewRoom(createRoom)
	if err.IsError() {
		log.Println(err.Error())
		return
	}
	room := model.NewRoom(createRoom.Name, createRoom.Private, client)

	// Get each client for members, doing the same in private chat when the member is offline
	// Client
	room.AddClients(p.clientManager.GetClientsByUserId(client.Id)...)
	// MemberIds
	for _, m := range createRoom.MemberIds {
		clients := p.clientManager.GetClientsByUserId(m)
		room.AddClients(clients...)
	}
	p.roomManager.AddRooms(room)

	// respond with room_id
	payload := model.NewPayload(model.PayloadCreateRoom, &outcome)
	client.SendPayload(&payload)
}

func (p *PayloadHandler) HandleJoinRoom(sender *model.Client, joinRoom *model.IncomeJoinRoom) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		return
	}
	if room.Private {
		return
	}

	if cerr := p.chatService.HandleJoinRoom(joinRoom); cerr.IsError() {
		return
	}

	room.AddClients(p.clientManager.GetClientsByUserId(sender.UserId)...)
	// Give notification to all users in room
	incomeNotif := model.IncomeNotification{
		Type:     model.NotifJoinRoom,
		Receiver: joinRoom.RoomId,
	}
	p.HandleRoomNotification(sender, &incomeNotif)
}

func (p *PayloadHandler) HandleLeaveRoom(sender *model.Client, joinRoom *model.IncomeJoinRoom) {
	room, err := p.roomManager.GetRoomById(joinRoom.RoomId)
	if err != nil {
		log.Println(err)
		return
	}
	room.RemoveClientsByUserId(sender.UserId)

	incomeNotif := model.IncomeNotification{
		Type:     model.NotifLeaveRoom,
		Receiver: joinRoom.RoomId,
	}

	p.HandleRoomNotification(sender, &incomeNotif)
}

func (p *PayloadHandler) HandleGetUsers(client *model.Client, userPayload *model.IncomeGetUser) {
	clients := p.clientManager.GetClientsByUsername(userPayload.Username)
	// Send back with the users
	outcomeGetUser := model.NewOutcomeGetUser(clients)
	payload := model.NewPayload(model.PayloadGetUsers, &outcomeGetUser)
	client.SendPayload(&payload)
}
