package handler

import (
	"log"

	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/util"
	"chatto/internal/util/containers"
	"chatto/internal/ws/manager"

	"chatto/internal/model"
	"chatto/internal/service"
)

func StartPayloadHandler(payload <-chan *model.Payload, chatService service.IChatService, roomManager *manager.RoomManager, clientManager *manager.ClientManager) {
	handler := &PayloadHandler{
		payload:       payload,
		roomManager:   roomManager,
		clientManager: clientManager,
		chatService:   chatService,
	}
	go handler.processPayload()
}

type PayloadHandler struct {
	payload <-chan *model.Payload

	// manager should not modify, it is supposed to only read the clients or rooms
	roomManager   *manager.RoomManager
	clientManager *manager.ClientManager

	chatService service.IChatService
}

func (p *PayloadHandler) processPayload() {
	for payload := range p.payload {
		switch payload.Type {
		case model.PayloadTyping:
			input, err := model.PayloadData[dto.TypingInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleTyping(payload.Sender, input)
		case model.PayloadMessage:
			roomChat, err := model.PayloadData[dto.MessageInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleRoomMessage(payload.Sender, &roomChat)
		case model.PayloadCreateRoom:
			createRoom, err := model.PayloadData[dto.CreateRoomInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleCreateRoom(payload.Sender, &createRoom)
		case model.PayloadJoinRoom:
			joinRoom, err := model.PayloadData[dto.RoomInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleJoinRoom(payload.Sender, joinRoom)
		case model.PayloadLeaveRoom:
			// Implicitly remove room when there is only one user there
			leaveRoom, err := model.PayloadData[dto.RoomInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleLeaveRoom(payload.Sender, leaveRoom)
		case model.PayloadInviteToRoom:
			inviteRoom, err := model.PayloadData[dto.MemberRoomInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleInviteToRoom(payload.Sender, inviteRoom)
		case model.PayloadKickFromRoom:
			kickRoom, err := model.PayloadData[dto.MemberRoomInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleKickFromRoom(payload.Sender, kickRoom)
		case model.PayloadGetUsers:
			getUser, err := model.PayloadData[dto.GetUserInput](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleGetUsers(payload.Sender, &getUser)
		case model.PayloadGetChats:
			getChats, err := model.PayloadData[dto.MessageRequest](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleMessageRequest(payload.Sender, &getChats)
		case model.PayloadGetNotifications:
			getNotifs, err := model.PayloadData[dto.NotificationRequest](payload)
			if err != nil {
				log.Println(err)
				continue
			}
			p.HandleGetNotifications(payload.Sender, &getNotifs)
		case model.PayloadGetUserRooms:
			p.HandleGetUserRooms(payload.Sender)
		default:
			output := model.NewErrorPayloadOutput(common.PAYLOAD_BAD_FORMAT_ERROR, constant.MSG_BAD_FORMAT_PAYLOAD)
			payload.Sender.SendPayload(&output)
		}
	}
	log.Println("Closed")
}

func handleForward[T any](p *PayloadHandler, sender *model.Client, types string, data *T) {
	senders := p.clientManager.GetUniqueClientByUserId(sender)
	if !containers.IsEmpty(senders) {
		forwardOutput := model.NewPayloadOutput(types, data)
		p.broadcast(&forwardOutput, senders)
	}
}

// broadcast Used for send message into all the receivers
func (p *PayloadHandler) broadcast(output *model.PayloadOutput, receivers []*model.Client) {
	for _, c := range receivers {
		c.SendPayload(output)
	}
}

func (p *PayloadHandler) handleNotification(senderUserId string, input *dto.NotificationInput) {
	notifOutput, cerr := p.chatService.NewNotification(senderUserId, input)
	if cerr.IsError() {
		log.Println(cerr.Error())
		return
	}

	// Broadcast
	room, _ := p.roomManager.GetRoomById(input.ReceiverId)
	payload := model.NewPayloadOutput(model.PayloadNotification, &notifOutput)
	room.Broadcast(&payload)
}

func (p *PayloadHandler) HandleTyping(sender *model.Client, input dto.TypingInput) {
	notifInput := dto.NotificationInput{
		Type:       model.NotifTyping,
		ReceiverId: input.RoomId,
	}
	p.handleNotification(sender.UserId, &notifInput)
}

func (p *PayloadHandler) HandleRoomMessage(sender *model.Client, input *dto.MessageInput) {
	messageOutput, cerr := p.chatService.NewMessage(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}

	// Broadcast
	room, _ := p.roomManager.GetRoomById(input.ReceiverId)
	payload := model.NewPayloadOutput(model.PayloadMessage, &messageOutput)
	room.Broadcast(&payload)
}

func (p *PayloadHandler) HandleCreateRoom(sender *model.Client, input *dto.CreateRoomInput) {
	output, err := p.chatService.CreateRoom(sender, input)
	if err.IsError() {
		util.SendErrorPayload(sender, err)
		return
	}

	// respond with room_id
	util.SendSuccessPayload(sender, &output)
}

func (p *PayloadHandler) HandleJoinRoom(sender *model.Client, input dto.RoomInput) {
	cerr := p.chatService.JoinRoom(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}

	// Give notification to all members that there are users joined into rooms
	notifInput := dto.NotificationInput{
		Type:       model.NotifJoinRoom,
		ReceiverId: input.RoomId,
	}
	p.handleNotification(sender.UserId, &notifInput)

	util.SendNilSuccessPayload(sender)
}

func (p *PayloadHandler) HandleLeaveRoom(sender *model.Client, input dto.RoomInput) {
	cerr := p.chatService.LeaveRoom(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}

	//Give notification to all members that this user is leaving
	notifInput := dto.NotificationInput{
		Type:       model.NotifLeaveRoom,
		ReceiverId: input.RoomId,
	}
	p.handleNotification(sender.UserId, &notifInput)
	util.SendNilSuccessPayload(sender)
}

func (p *PayloadHandler) HandleInviteToRoom(sender *model.Client, input dto.MemberRoomInput) {
	if input.UserIds == nil || containers.IsEmpty(input.UserIds) {
		util.SendErrorPayload(sender, common.NewError(common.PAYLOAD_BAD_FORMAT_ERROR, constant.MSG_BAD_FORMAT_PAYLOAD))
		return
	}

	cerr := p.chatService.Invite(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}

	for _, userId := range input.UserIds {
		// Give notification to all users in room
		notifInput := dto.NotificationInput{
			Type:       model.NotifJoinRoom,
			ReceiverId: input.RoomId,
		}

		p.handleNotification(userId, &notifInput)
	}

	util.SendNilSuccessPayload(sender)
}

func (p *PayloadHandler) HandleKickFromRoom(sender *model.Client, input dto.MemberRoomInput) {
	if input.UserIds == nil || containers.IsEmpty(input.UserIds) {
		util.SendErrorPayload(sender, common.NewError(common.PAYLOAD_BAD_FORMAT_ERROR, constant.MSG_BAD_FORMAT_PAYLOAD))
		return
	}

	cerr := p.chatService.KickOut(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}

	for _, userId := range input.UserIds {
		// Give notification to all users in room
		notifInput := dto.NotificationInput{
			Type:       model.NotifLeaveRoom,
			ReceiverId: input.RoomId,
		}

		p.handleNotification(userId, &notifInput)
	}
	util.SendNilSuccessPayload(sender)
}

func (p *PayloadHandler) HandleGetUsers(sender *model.Client, input *dto.GetUserInput) {
	// Send back with the users
	output, cerr := p.chatService.GetUsersByName(sender, input.Username)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
	} else {
		util.SendSuccessPayload(sender, &output)
	}
}

func (p *PayloadHandler) HandleMessageRequest(sender *model.Client, input *dto.MessageRequest) {
	messages, cerr := p.chatService.GetRoomMessages(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}
	util.SendSuccessPayload(sender, &messages)
}

func (p *PayloadHandler) HandleGetNotifications(sender *model.Client, input *dto.NotificationRequest) {
	notifs, cerr := p.chatService.GetRoomNotifications(sender, input)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}
	util.SendSuccessPayload(sender, &notifs)
}

func (p *PayloadHandler) HandleGetUserRooms(sender *model.Client) {
	rooms, cerr := p.chatService.GetRoomsByUserId(sender.UserId)
	if cerr.IsError() {
		util.SendErrorPayload(sender, cerr)
		return
	}
	util.SendSuccessPayload(sender, &rooms)
}
