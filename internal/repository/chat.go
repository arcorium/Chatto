package repository

import (
	"server_client_chat/internal/model"

	"github.com/redis/go-redis/v9"
)

func NewChatRepository(client *redis.Client) ChatRepository {
	return ChatRepository{db: client}
}

type ChatRepository struct {
	db *redis.Client
}

func (c ChatRepository) InsertChat(message *model.Message) error {
	return nil
}

func (c ChatRepository) FindChats() ([]model.Message, error) {
	return nil, nil
}
