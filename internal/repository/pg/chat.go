package pg_repo

import (
	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/repository"
	"chatto/internal/util"

	"github.com/redis/go-redis/v9"
)

func NewChatRepository(client *redis.Client) repository.IChatRepository {
	return &chatRepository{db_: client}
}

type chatRepository struct {
	db_ *redis.Client
}

func (c chatRepository) db() *redis.Client {
	return c.db_
}

func (c chatRepository) NewClient(client *model.Client) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()
	result := c.db().SAdd(ctx, constant.REDIS_KEY_USER, *client)
	return result.Err()
}

func (c chatRepository) RemoveClient(client *model.Client) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()
	result := c.db().SRem(ctx, constant.REDIS_KEY_USER, *client)
	return result.Err()
}

func (c chatRepository) CreateMessage(id string, message *model.Message) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	result := c.db().SAdd(ctx, "chat:"+id, *message)
	return result.Err()
}

func (c chatRepository) CreateNotification(id string, notif *model.Notification) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	result := c.db().SAdd(ctx, "notif:"+id, *notif)
	return result.Err()
}

func (c chatRepository) FindChats() ([]model.Message, error) {
	return nil, nil
}
