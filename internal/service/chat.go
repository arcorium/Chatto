package service

import (
	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/repository"
)

type IChatService interface {
	HandleNewClient(client *model.Client) common.Error
	HandleNewPrivateMessage(sender *model.Client, message *model.IncomeMessage) (model.OutcomePrivateMessage, common.Error)
	HandleNewRoomMessage(sender *model.Client, message *model.IncomeMessage) (model.OutcomeRoomMessage, common.Error)
	HandleNewPrivateNotification(sender *model.Client, message *model.IncomeNotification) (model.OutcomePrivateNotification, common.Error)
	HandleNewRoomNotification(sender *model.Client, message *model.IncomeNotification) (model.OutcomeRoomNotification, common.Error)
	HandleNewRoom(room *model.IncomeCreateRoom) (model.OutcomeCreateRoom, common.Error)
	HandleInviteToRoom(room *model.IncomeInviteRoom) common.Error
	HandleJoinRoom(room *model.IncomeJoinRoom) common.Error
	HandleLeaveRoom(room *model.IncomeJoinRoom) common.Error
}

func NewChatService(chatRepository repository.ChatRepository) ChatService {
	return ChatService{chatRepo: chatRepository}
}

type ChatService struct {
	chatRepo repository.ChatRepository
}

func (c ChatService) HandleNewPrivateNotification(sender *model.Client, message *model.IncomeNotification) (model.OutcomePrivateNotification, common.Error) {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleNewRoomNotification(sender *model.Client, message *model.IncomeNotification) (model.OutcomeRoomNotification, common.Error) {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleNewRoom(room *model.IncomeCreateRoom) (model.OutcomeCreateRoom, common.Error) {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleInviteToRoom(room *model.IncomeInviteRoom) common.Error {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleJoinRoom(room *model.IncomeJoinRoom) common.Error {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleLeaveRoom(room *model.IncomeJoinRoom) common.Error {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleNewClient(client *model.Client) common.Error {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) HandleNewPrivateMessage(sender *model.Client, message *model.IncomeMessage) (model.OutcomePrivateMessage, common.Error) {
	// message for storing into database
	chatMessage := model.NewMessage(sender, message)
	if err := c.chatRepo.UpsertMessage(chatMessage, false); err != nil {
		return model.OutcomePrivateMessage{}, common.NewError(constant.INTERNAL_ERROR, err.Error())
	}

	return model.NewOutcomePrivateMessage(sender, chatMessage), common.NoError()
}

func (c ChatService) HandleNewRoomMessage(sender *model.Client, message *model.IncomeMessage) (model.OutcomeRoomMessage, common.Error) {
	chatMessage := model.NewMessage(sender, message)
	if err := c.chatRepo.UpsertMessage(chatMessage, true); err != nil {
		return model.OutcomeRoomMessage{}, common.NewError(constant.INTERNAL_ERROR, err.Error())
	}

	return model.NewOutcomeRoomMessage(sender, chatMessage), common.NoError()
}
