package service

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/repository"
	"chatto/internal/util/ctrutil"
)

type IChatService interface {
	NewClient(client *model.Client) common.Error
	RemoveClient(client *model.Client) common.Error
	NewMessage(sender *model.Client, message *dto.MessageInput) (dto.MessageOutput, common.Error)
	NewNotification(senderId string, input *dto.NotificationInput) (dto.NotificationOutput, common.Error)
	GetUsersByName(name string) (dto.GetUserOutput, common.Error)
}

func NewChatService(chatRepository repository.IChatRepository) ChatService {
	return ChatService{chatRepo: chatRepository}
}

type ChatService struct {
	chatRepo repository.IChatRepository
	userRepo repository.IUserRepository
}

func (c ChatService) GetUsersByName(name string) (dto.GetUserOutput, common.Error) {
	// TODO: Get clients on redis database to check whether the user is online or not
	// Get all from userRepo
	users, err := c.userRepo.FindUsersByLikelyName(name)
	if err != nil {
		return dto.GetUserOutput{}, common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	userResponses := ctrutil.ConvertSliceType(users, func(current *model.User) dto.UserResponse {
		return dto.NewUserResponse(current)
	})

	return dto.NewGetUserOutput(userResponses), common.NoError()
}

func (c ChatService) NewNotification(senderId string, input *dto.NotificationInput) (dto.NotificationOutput, common.Error) {
	notif := dto.NewNotificationFromInput(senderId, input)

	if notif.Type == model.NotifJoinRoom || notif.Type == model.NotifLeaveRoom {
		_ = c.chatRepo.CreateNotification(input.Receiver, &notif)
	}

	output := dto.NewNotificationOutput(senderId, input, &notif)
	return output, common.NoError()
}

func (c ChatService) NewClient(client *model.Client) common.Error {
	err := c.chatRepo.NewClient(client)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (c ChatService) RemoveClient(client *model.Client) common.Error {
	err := c.chatRepo.RemoveClient(client)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (c ChatService) NewMessage(sender *model.Client, input *dto.MessageInput) (dto.MessageOutput, common.Error) {
	// message for storing into database
	message := dto.NewMessageFromInput(sender, input)
	if err := c.chatRepo.CreateMessage(input.Receiver, &message); err != nil {
		return dto.MessageOutput{}, common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	return dto.NewMessageOutput(input, &message), common.NoError()
}
