package service

import (
	"log"

	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/repository"
	"chatto/internal/util/containers"
	"chatto/internal/ws/manager"
)

type IChatService interface {
	// ProcessPayload Used for processing raw message, like decryption and providing information, so it can processed further
	ProcessPayload(sender *model.Client, input *model.PayloadInput) model.Payload
	// NewClient Used for when new sender online, it will interact with clientManager and roomManager
	NewClient(sender *model.Client) common.Error
	CreateRoom(sender *model.Client, input *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error)
	RemoveClient(sender *model.Client) common.Error
	NewMessage(sender *model.Client, message *dto.MessageInput) (dto.MessageOutput, common.Error)
	NewNotification(senderUserId string, input *dto.NotificationInput) (dto.NotificationOutput, common.Error)
	GetUsersByName(sender *model.Client, name string) (dto.GetUserOutput, common.Error)
	GetRoomsByUserId(userId string) ([]dto.UserRoomResponse, common.Error)
	GetRoomMessages(sender *model.Client, request *dto.MessageRequest) ([]dto.MessageResponse, common.Error)
	GetRoomNotifications(sender *model.Client, request *dto.NotificationRequest) ([]dto.NotificationResponse, common.Error)
	JoinRoom(sender *model.Client, input dto.RoomInput) common.Error
	LeaveRoom(sender *model.Client, input dto.RoomInput) common.Error
	KickOut(sender *model.Client, input dto.MemberRoomInput) common.Error
	Invite(sender *model.Client, input dto.MemberRoomInput) common.Error
	ClearUsers() common.Error
}

func NewChatService(chatRepository repository.IChatRepository, userService IUserService, roomService IRoomService, roomManager *manager.RoomManager, clientManager *manager.ClientManager) IChatService {
	return &chatService{
		repo:          chatRepository,
		roomManager:   roomManager,
		clientManager: clientManager,
		userService:   userService,
		roomService:   roomService,
	}
}

type chatService struct {
	repo repository.IChatRepository

	roomManager   *manager.RoomManager
	clientManager *manager.ClientManager

	userService IUserService
	roomService IRoomService
}

func (c *chatService) ProcessPayload(sender *model.Client, input *model.PayloadInput) model.Payload {
	return model.Payload{
		Type:   input.Type,
		Data:   input.Data,
		Sender: sender,
	}
}

func (c *chatService) CreateRoom(sender *model.Client, input *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error) {
	// Prevent the sender as memberIds
	if containers.SliceContains(input.MemberIds, sender.UserId) {
		return dto.CreateRoomOutput{}, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	// Create room on room service
	output, cerr := c.roomService.CreateRoom(input)
	if cerr.IsError() {
		return dto.CreateRoomOutput{}, cerr
	}

	// Add creator as admin
	role := model.RoomRoleAdmin
	if input.Private {
		role = model.RoomRoleUser
	}
	userRoomInput := dto.UserRoomAddInput{
		Users:  dto.NewUserRoles(role, sender.UserId),
		RoomId: output.Id,
	}
	c.roomService.AddUsersInRoom(userRoomInput, true)

	// Handle for room manager
	creator := c.clientManager.GetClientsByUserId(sender.UserId)

	members := make([]*model.Client, 0, len(input.MemberIds)*2)
	for _, userId := range input.MemberIds {
		clients := c.clientManager.GetClientsByUserId(userId)
		members = append(clients)
	}
	chatRoom := dto.NewChatRoomFromOutput(&output, members...)
	// Add creator as admin
	chatRoom.AddClientsWithSameRole(role, creator...)

	c.roomManager.AddRooms(chatRoom)
	return output, common.NoError()
}

func (c *chatService) JoinRoom(sender *model.Client, input dto.RoomInput) common.Error {
	// Check room existences
	room, err := c.roomManager.GetRoomById(input.RoomId)
	if err != nil {
		return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	// Check if user already room's member
	if room.IsClientExist(sender) {
		return common.NewError(common.USER_ALREADY_ROOM_MEMBER, constant.MSG_USER_ALREADY_ROOM_MEMBER)
	}

	if room.Private {
		return common.NewError(common.ROOM_IS_PRIVATE_ERROR, constant.MSG_JOIN_PRIVATE_ROOM)
	}

	userRoomInput := dto.UserRoomAddInput{
		Users:  dto.NewUserRoles(model.RoomRoleUser, sender.UserId),
		RoomId: input.RoomId,
	}

	cerr := c.roomService.AddUsersInRoom(userRoomInput, false)
	if cerr.IsError() {
		return cerr
	}

	clients := c.clientManager.GetClientsByUserId(sender.UserId)
	room.AddClientsWithSameRole(model.RoomRoleUser, clients...)

	return common.NoError()
}

func (c *chatService) LeaveRoom(sender *model.Client, input dto.RoomInput) common.Error {
	room, err := c.roomManager.GetRoomById(input.RoomId)
	if err != nil {
		return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	if !room.IsClientExist(sender) {
		return common.NewError(common.USER_NOT_ROOM_MEMBER, constant.MSG_USER_NOT_ROOM_MEMBER)
	}

	if room.Private {
		return common.NewError(common.ROOM_IS_PRIVATE_ERROR, constant.MSG_LEAVE_FROM_PRIVATE_ROOM)
	}

	userRoomInput := dto.UserRoomRemoveInput{
		UserIds: []string{sender.UserId},
		RoomId:  input.RoomId,
	}
	cerr := c.roomService.RemoveUsersInRoom(userRoomInput)
	if cerr.IsError() {
		return cerr
	}

	// Remove from room manager
	room.RemoveClientsByUserId(sender.UserId)

	if containers.IsEmpty(room.UserIds()) {
		err = c.roomManager.RemoveRoomById(room.Id)
	}

	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (c *chatService) KickOut(sender *model.Client, input dto.MemberRoomInput) common.Error {
	// Prevent the sender in user_ids
	if containers.SliceContains(input.UserIds, sender.UserId) {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	room, cerr := c.isAdmin(sender, input.RoomId)
	if cerr.IsError() {
		return cerr
	}

	if room.Private {
		return common.NewError(common.ROOM_IS_PRIVATE_ERROR, constant.MSG_KICK_FROM_PRIVATE_ROOM)
	}

	users, cerr := c.roomService.FindRoomMembersById(input.RoomId)
	if cerr.IsError() {
		return cerr
	}

	// Check if users is there
	userIds := containers.ConvertSlice(users, func(current *dto.UserResponse) string {
		return current.Id
	})
	if !containers.SliceContains(userIds, input.UserIds...) {
		return common.NewError(common.MEMBER_ROOM_NOT_FOUND_ERROR, constant.MSG_MEMBER_ROOM_NOT_FOUND)
	}

	userRoomInput := dto.UserRoomRemoveInput{
		RoomId:  input.RoomId,
		UserIds: input.UserIds,
	}
	cerr = c.roomService.RemoveUsersInRoom(userRoomInput)
	if cerr.IsError() {
		return cerr
	}

	// Handle manager
	for _, userId := range input.UserIds {
		room.RemoveClientsByUserId(userId)
	}

	return common.NoError()
}

func (c *chatService) Invite(sender *model.Client, input dto.MemberRoomInput) common.Error {
	// Prevent the sender in user_ids
	if containers.SliceContains(input.UserIds, sender.UserId) {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	room, cerr := c.isAdmin(sender, input.RoomId)
	if cerr.IsError() {
		return cerr
	}

	if room.Private {
		return common.NewError(common.ROOM_IS_PRIVATE_ERROR, constant.MSG_INVITE_TO_PRIVATE_ROOM)
	}

	userRoomInput := dto.UserRoomAddInput{
		Users:  dto.NewUserRoles(model.RoomRoleUser, input.UserIds...),
		RoomId: input.RoomId,
	}

	// Duplication handled by the repository
	cerr = c.roomService.AddUsersInRoom(userRoomInput, true)
	if cerr.IsError() {
		return cerr
	}

	// Handle room manager
	for _, userId := range input.UserIds {
		clients := c.clientManager.GetClientsByUserId(userId)
		room.AddClientsWithSameRole(model.RoomRoleUser, clients...)
	}
	return common.NoError()
}

func (c *chatService) GetRoomsByUserId(userId string) ([]dto.UserRoomResponse, common.Error) {
	return c.roomService.FindUserRoomsByUserId(userId)
}

func (c *chatService) GetRoomMessages(sender *model.Client, request *dto.MessageRequest) ([]dto.MessageResponse, common.Error) {
	cerr := c.checkRoomAndUserExistences(sender, request.RoomId)
	if cerr.IsError() {
		return nil, cerr
	}

	messages, err := c.repo.FindRoomChats(request)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	messageResponses := containers.ConvertSlice(messages, dto.NewMessageResponse)
	return messageResponses, common.NoError()
}

func (c *chatService) GetRoomNotifications(sender *model.Client, request *dto.NotificationRequest) ([]dto.NotificationResponse, common.Error) {
	cerr := c.checkRoomAndUserExistences(sender, request.RoomId)
	if cerr.IsError() {
		return nil, cerr
	}

	notifs, err := c.repo.FindRoomNotifications(request)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	notifsResponses := containers.ConvertSlice(notifs, dto.NewNotificationResponse)
	return notifsResponses, common.NoError()
}

func (c *chatService) GetUsersByName(sender *model.Client, name string) (dto.GetUserOutput, common.Error) {
	// Get all from userService
	userResponses, err := c.userService.FindUsersByLikelyName(name)
	if err.IsError() {
		return dto.GetUserOutput{}, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	// Filter to prevent it returns back the sender name
	userResponses = containers.SliceFilter(userResponses, func(current *dto.UserResponse) bool {
		return current.Id != sender.Id
	})

	return dto.NewGetUserOutput(userResponses), common.NoError()
}

func (c *chatService) NewNotification(senderUserId string, input *dto.NotificationInput) (dto.NotificationOutput, common.Error) {
	// Check room existences
	_, err := c.roomManager.GetRoomById(input.ReceiverId)
	if err != nil {
		return dto.NotificationOutput{}, common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	notif := dto.NewNotificationFromInput(senderUserId, input)
	if notif.Type == model.NotifJoinRoom || notif.Type == model.NotifLeaveRoom {
		_ = c.repo.CreateNotification(&notif)
	}

	output := dto.NewNotificationOutput(&notif)
	return output, common.NoError()
}

func (c *chatService) NewClient(sender *model.Client) common.Error {
	roomResponses, cerr := c.roomService.FindUserRoomsByUserId(sender.UserId)
	if cerr.IsError() {
		return cerr
	}

	for _, roomResponse := range roomResponses {
		chatRoom, err := c.roomManager.GetRoomById(roomResponse.RoomId)
		if err != nil {
			log.Println(err)
			return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		}
		chatRoom.AddClient(sender, roomResponse.UserRole)
	}

	c.clientManager.AddClients(sender)

	// Send list of sender rooms
	output := model.NewPayloadOutput(model.PayloadGetUserRooms, &roomResponses)
	sender.SendPayload(&output)

	err := c.repo.NewClient(sender)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (c *chatService) RemoveClient(sender *model.Client) common.Error {
	roomResponses, cerr := c.roomService.FindUserRoomsByUserId(sender.UserId)
	if cerr.IsError() {
		return cerr
	}

	for _, roomResponse := range roomResponses {
		chatRoom, err := c.roomManager.GetRoomById(roomResponse.RoomId)
		if err != nil {
			log.Println(err)
			return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		}
		chatRoom.RemoveClient(sender)
	}

	c.clientManager.RemoveClientById(sender.Id)

	err := c.repo.RemoveClient(sender)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (c *chatService) NewMessage(sender *model.Client, input *dto.MessageInput) (dto.MessageOutput, common.Error) {
	cerr := c.checkRoomAndUserExistences(sender, input.ReceiverId)
	if cerr.IsError() {
		return dto.MessageOutput{}, cerr
	}

	// message for storing into database
	message := dto.NewMessageFromInput(sender, input)
	if err := c.repo.CreateMessage(&message); err != nil {
		return dto.MessageOutput{}, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	return dto.NewMessageOutput(&message), common.NoError()
}

func (c *chatService) ClearUsers() common.Error {
	err := c.repo.ResetClients()
	if err != nil {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	return common.NoError()
}

// checkRoomAndUserExistences Used to check if the room is exists and if the sender is the room's member
func (c *chatService) checkRoomAndUserExistences(sender *model.Client, roomId string) common.Error {
	room, err := c.roomManager.GetRoomById(roomId)
	if err != nil {
		return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	if !room.IsClientExist(sender) {
		return common.NewError(common.USER_NOT_ROOM_MEMBER, constant.MSG_USER_NOT_ROOM_MEMBER)
	}
	return common.NoError()
}

// isAdmin Used to check if the sender is Admin on the room, it will return the ChatRoom and error based on where the error occurred, the error can occure on
func (c *chatService) isAdmin(sender *model.Client, roomId string) (*model.ChatRoom, common.Error) {
	// Get room existence
	room, err := c.roomManager.GetRoomById(roomId)
	if err != nil {
		return nil, common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	// Get sender role
	role, exist := room.GetRoleByUserId(sender.UserId)
	if !exist {
		return nil, common.NewError(common.USER_NOT_ROOM_MEMBER, constant.MSG_USER_NOT_ROOM_MEMBER)
	}

	// Check role
	if role != model.RoomRoleAdmin {
		return nil, common.NewError(common.AUTH_UNAUTHORIZED, constant.MSG_AUTH_UNAUTHORIZED)
	}

	return room, common.NoError()
}
