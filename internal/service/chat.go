package service

import (
	"server_client_chat/internal/model"
	"server_client_chat/internal/repository"
)

func NewChatService(chatRepository repository.ChatRepository) ChatService {
	return ChatService{chatRepo: chatRepository}
}

type ChatService struct {
	chatRepo repository.ChatRepository
}

func (c ChatService) HandleNewClient(client *model.Client) CustomError {
	//TODO implement me
	panic("implement me")
}
