package handler

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/util"
	"log"

	"chatto/internal/model"
	"chatto/internal/service"
	"chatto/internal/ws/manager"
)

func StartPayloadHandler(clientManager *manager.ClientManager, roomManager *manager.RoomManager, payload <-chan *model.Payload, chatService service.IChatService, roomService service.IRoomService) {
	handler := &PayloadHandler{
		roomManager:   roomManager,
		clientManager: clientManager,
		payload:       payload,
		chatService:   chatService,
		roomService:   roomService,
	}

	go handler.PayloadHandler()
}

type PayloadHandler struct {
	roomManager   *manager.RoomManager
	clientManager *manager.ClientManager
	payload       <-chan *model.Payload

	chatService service.IChatService
	roomService service.IRoomService
}

func (p *PayloadHandler) PayloadHandler() {
	for {
		select {
		case payload := <-p.payload:
			log.Println(payload)

			// Convert interface{} into bytes, so it can be unmarshalled
			bytes, err := payload.DataBytes()
			if err != nil {
				log.Println(err)
				continue
			}

			// Check payload type
			switch payload.Type {
			case model.PayloadPrivateNotification:
				notif, err := model.Decode[dto.NotificationInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandlePrivateNotification(payload.Client, &notif)
			case model.PayloadRoomNotification:
				notif, err := model.Decode[dto.NotificationInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleRoomNotification(payload.Client, &notif)
			case model.PayloadPrivateChat:
				privateChat, err := model.Decode[dto.MessageInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandlePrivateChat(payload.Client, &privateChat)
				// Let client handle it
			case model.PayloadRoomChat:
				roomChat, err := model.Decode[dto.MessageInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleRoomChat(payload.Client, &roomChat)
			case model.PayloadCreateRoom:
				createRoom, err := model.Decode[dto.CreateRoomInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleCreateRoom(payload.Client, &createRoom)
			case model.PayloadJoinRoom:
				joinRoom, err := model.Decode[dto.JoinRoomInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleJoinRoom(payload.Client, &joinRoom)
			case model.PayloadLeaveRoom:
				// Implicitly remove room when there is only one user there
				leaveRoom, err := model.Decode[dto.LeaveRoomInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleLeaveRoom(payload.Client, &leaveRoom)
			case model.PayloadInviteToRoom:
				inviteRoom, err := model.Decode[dto.InviteRoomInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleInviteToRoom(payload.Client, &inviteRoom)
			case model.PayloadGetUsers:
				getUser, err := model.Decode[dto.GetUserInput](bytes)
				if err != nil {
					log.Println(err)
					continue
				}
				p.HandleGetUsers(payload.Client, &getUser)
			}
		}
	}
}

func (p *PayloadHandler) broadcast(output *model.PayloadOutput, receivers []*model.Client) {
	for _, c := range receivers {
		c.SendPayload(output)
	}
}

func (p *PayloadHandler) HandlePrivateNotification(sender *model.Client, input *dto.NotificationInput) {
	receivers := p.clientManager.GetClientsByUserId(input.Receiver)
	if len(receivers) == 0 {
		return
	}

	privateNotif, err := p.chatService.NewNotification(sender.UserId, input)
	if err.IsError() {
		return
	}

	senders := p.clientManager.GetUniqueClientByUserId(sender)
	if len(senders) == 0 {
		forwardNotification := dto.NewNotificationForward(input.Receiver, &privateNotif)
		util.ForwardPayload(senders, model.PayloadNotificationForwarder, &forwardNotification)
	}

	notifPayload := model.NewPayloadOutput(model.PayloadPrivateChat, &privateNotif)
	p.broadcast(&notifPayload, receivers)
}

func (p *PayloadHandler) HandleRoomNotification(sender *model.Client, input *dto.NotificationInput) {
	// Search rooms
	room, err := p.roomManager.GetRoomById(input.Receiver)
	if err != nil {
		return
	}

	notifOutput, err := p.chatService.NewNotification(sender.UserId, input)
	if err != nil {
		return
	}

	// Broadcast
	payload := model.NewPayloadOutput(model.PayloadRoomNotification, &notifOutput)
	room.Broadcast(&payload, sender.Id)
}

func (p *PayloadHandler) HandlePrivateChat(sender *model.Client, input *dto.MessageInput) {
	receivers := p.clientManager.GetClientsByUserId(input.Receiver)

	privateMessage, err := p.chatService.NewMessage(sender, input)
	if err.IsError() {
		return
	}

	// Send to all clients for the same user id
	senders := p.clientManager.GetUniqueClientByUserId(sender)
	if len(senders) >= 1 {
		forwardMessage := dto.NewMessageForward(input.Receiver, &privateMessage)
		util.ForwardPayload(senders, model.PayloadMessageForwarder, &forwardMessage)
	}

	payloadOutput := model.NewPayloadOutput(model.PayloadPrivateChat, &privateMessage)
	p.broadcast(&payloadOutput, receivers)
}

func (p *PayloadHandler) HandleRoomChat(sender *model.Client, input *dto.MessageInput) {
	// Search rooms
	room, _ := p.roomManager.GetRoomById(input.Receiver)

	messageOutput, cerr := p.chatService.NewMessage(sender, input)
	if cerr.IsError() {
		log.Println(cerr.Error())
	}

	// Broadcast
	payload := model.NewPayloadOutput(model.PayloadRoomChat, &messageOutput)
	room.Broadcast(&payload, sender.Id)
}

func (p *PayloadHandler) HandleCreateRoom(sender *model.Client, input *dto.CreateRoomInput) {
	output, err := p.roomService.CreateRoom(input)
	if err.IsError() {
		log.Println(err.Error())
		return
	}

	// Get all clients for sender and the members
	members := make([]*model.Client, 0, len(input.MemberIds)*2)
	senders := p.clientManager.GetClientsByUserId(sender.Id)
	members = append(senders)
	for _, m := range input.MemberIds {
		clients := p.clientManager.GetClientsByUserId(m)
		members = append(clients)
	}

	room := dto.NewChatRoom(&output, members...)
	p.roomManager.AddRooms(room)

	// respond with room_id
	util.SendSuccessPayload(sender, &output)
}

func (p *PayloadHandler) HandleJoinRoom(sender *model.Client, input *dto.JoinRoomInput) {
	room, err := p.roomManager.GetRoomById(input.RoomId)
	if err != nil {
		util.SendErrorPayload(sender, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		return
	}

	// Check if already joined
	// TODO: Implement it on http too
	if room.IsUserExist(sender.UserId) {
		room.AddClients(sender) // Add clients already handle multiple same client perfectly
		util.SendErrorPayload(sender, common.USER_ALREADY_ROOM_MEMBER, constant.MSG_JOIN_JOINED_ROOM)
		return
	}

	if room.Private {
		util.SendErrorPayload(sender, common.ROOM_IS_PRIVATE_ERROR, constant.MSG_JOIN_PRIVATE_ROOM)
		return
	}

	userIds, _ := p.roomService.AddUserIntoRoom(room.Id, dto.UserRoomInput{UserIds: []string{sender.UserId}})
	for _, userId := range userIds.UserIds {
		room.AddClients(p.clientManager.GetClientsByUserId(userId)...)
		// Give notification to all users in room
		notifInput := dto.NotificationInput{
			Type:     model.NotifJoinRoom,
			Receiver: input.RoomId,
		}

		notifOutput, _ := p.chatService.NewNotification(userId, &notifInput)

		// Broadcast
		payload := model.NewPayloadOutput(model.PayloadRoomNotification, &notifOutput)
		room.Broadcast(&payload, userId)
	}
}

func (p *PayloadHandler) HandleLeaveRoom(sender *model.Client, input *dto.LeaveRoomInput) {
	room, err := p.roomManager.GetRoomById(input.RoomId)
	if err != nil {
		util.SendErrorPayload(sender, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		return
	}
	// Check if room's member
	// TODO: Implement it on http too
	if !room.IsUserExist(sender.UserId) {
		util.SendErrorPayload(sender, common.USER_NOT_ROOM_MEMBER, constant.MSG_NOT_ROOM_MEMBER)
		return
	}

	p.roomService.RemoveUsersOnRoom(room.Id, dto.UserRoomInput{UserIds: []string{sender.UserId}})

	room.RemoveClientsByUserId(sender.UserId)

	notifInput := dto.NotificationInput{
		Type:     model.NotifLeaveRoom,
		Receiver: input.RoomId,
	}

	notifOutput, _ := p.chatService.NewNotification(sender.UserId, &notifInput)

	// Broadcast
	payload := model.NewPayloadOutput(model.PayloadRoomNotification, &notifOutput)
	room.Broadcast(&payload, sender.UserId)
}

func (p *PayloadHandler) HandleInviteToRoom(sender *model.Client, input *dto.InviteRoomInput) {
	room, err := p.roomManager.GetRoomById(input.RoomId)
	if err != nil {
		util.SendErrorPayload(sender, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		return
	}

	// Check if the sender is group member
	// TODO: Implement it on http too
	if !room.IsUserExist(sender.UserId) {
		util.SendErrorPayload(sender, common.USER_NOT_ROOM_MEMBER, constant.MSG_NOT_ROOM_MEMBER)
		return
	}

	userIds, err := p.roomService.AddUserIntoRoom(room.Id, dto.UserRoomInput{UserIds: []string{sender.UserId}})
	for _, userId := range userIds.UserIds {
		room.AddClients(p.clientManager.GetClientsByUserId(userId)...)
		// Give notification to all users in room
		notifInput := dto.NotificationInput{
			Type:     model.NotifJoinRoom,
			Receiver: input.RoomId,
		}

		notifOutput, _ := p.chatService.NewNotification(userId, &notifInput)

		// Broadcast
		payload := model.NewPayloadOutput(model.PayloadRoomNotification, &notifOutput)
		room.Broadcast(&payload, userId)
	}
}

func (p *PayloadHandler) HandleGetUsers(sender *model.Client, input *dto.GetUserInput) {
	// Send back with the users
	output, err := p.chatService.GetUsersByName(input.Username)
	if err.IsError() {
		util.SendErrorPayload(sender, common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
	} else {
		util.SendSuccessPayload(sender, &output)
	}
}
