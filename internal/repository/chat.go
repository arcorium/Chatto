package repository

import (
	"chatto/internal/model"

	"github.com/redis/go-redis/v9"
)

type IChatRepository interface {
	UpsertMessage(message *model.Message, isRoom bool) error
	FindChats() ([]model.Message, error)
}

func NewChatRepository(client *redis.Client) ChatRepository {
	return ChatRepository{db: client}
}

type ChatRepository struct {
	db *redis.Client
}

func (c ChatRepository) UpsertMessage(message *model.Message, isRoom bool) error {
	return nil
}

func (c ChatRepository) FindChats() ([]model.Message, error) {
	return nil, nil
}
