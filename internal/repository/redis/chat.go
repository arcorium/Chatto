package redis_repo

import (
	"encoding/json"
	"fmt"

	"chatto/internal/constant"
	"chatto/internal/dto"
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
	key := constant.REDIS_KEY_USER + client.UserId

	// Check existences
	resExist := c.db().Exists(ctx, key)
	if resExist.Err() != nil {
		return resExist.Err()
	}
	// Only increment online field when the key is exists
	if resExist.Val() != 0 {
		result := c.db().HIncrBy(ctx, key, "online", 1)
		if result.Err() != nil {
			return result.Err()
		}
		// Renew TTL
		resultBool := c.db().Expire(ctx, key, constant.USER_CHAT_EXPIRATION_DURATION)
		return resultBool.Err()
	}

	result := c.db().HMSet(ctx, key, "name", client.Username, "role", string(client.Role), "online", 1)
	if result.Err() != nil {
		return result.Err()
	}
	// Set TTL
	result = c.db().Expire(ctx, key, constant.USER_CHAT_EXPIRATION_DURATION)
	return result.Err()
}

func (c chatRepository) RemoveClient(client *model.Client) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()
	key := constant.REDIS_KEY_USER + client.UserId

	result := c.db().HIncrBy(ctx, key, "online", -1)
	return result.Err()
}

func (c chatRepository) CreateMessage(message *model.Message) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	key := constant.REDIS_KEY_CHAT + message.Id
	//message.SenderId = util.EscapeMinesSymbols(message.SenderId)
	//message.ReceiverId = util.EscapeMinesSymbols(message.ReceiverId)
	//message.Message = util.EscapeMinesSymbols(message.Message)
	bytes, err := json.Marshal(*message)
	if err != nil {
		return err
	}

	result := c.db().Do(ctx, "JSON.SET", key, "$", bytes)
	return result.Err()
}

func (c chatRepository) CreateNotification(notif *model.Notification) error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	key := constant.REDIS_KEY_NOTIF + notif.Id
	//notif.SenderId = util.EscapeMinesSymbols(notif.SenderId)
	//notif.ReceiverId = util.EscapeMinesSymbols(notif.ReceiverId)
	bytes, err := json.Marshal(*notif)
	if err != nil {
		return err
	}

	result := c.db().Do(ctx, "JSON.SET", key, "$", bytes)
	return result.Err()
}

func (c chatRepository) FindRoomChats(request *dto.MessageRequest) ([]model.Message, error) {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	query := fmt.Sprintf("@receiver:{%s}", util.EscapeMinesSymbols(request.RoomId))
	if request.ToTime == 0 {
		query += fmt.Sprintf(" @timestamp:[%d +inf]", request.FromTime)
	} else {
		query += fmt.Sprintf(" @timestamp:[%d %d]", request.FromTime, request.ToTime)
	}
	result := c.db().Do(ctx, "FT.SEARCH", constant.REDIS_KEY_CHAT_INDEX, query, "SORTBY", "timestamp", "ASC")
	if result.Err() != nil {
		return nil, result.Err()
	}
	_, messages := ftSearchConvert[model.Message](result.Val())

	return messages, nil
}

func (c chatRepository) FindRoomNotifications(request *dto.NotificationRequest) ([]model.Notification, error) {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	query := fmt.Sprintf("@receiver:{%s}", util.EscapeMinesSymbols(request.RoomId))
	if request.ToTime == 0 {
		query += fmt.Sprintf(" @timestamp:[%d +inf]", request.FromTime)
	} else {
		query += fmt.Sprintf(" @timestamp:[%d %d]", request.FromTime, request.ToTime)
	}
	result := c.db().Do(ctx, "FT.SEARCH", constant.REDIS_KEY_NOTIF_INDEX, query, "SORTBY", "timestamp", "ASC")
	if result.Err() != nil {
		return nil, result.Err()
	}

	_, notifs := ftSearchConvert[model.Notification](result.Val())

	return notifs, nil
}

func (c chatRepository) ResetClients() error {
	ctx, cancel := util.NewTimeoutContext()
	defer cancel()

	// Check existences
	var cursor uint64 = 0
	for {
		keys, newCursor, err := c.db().Scan(ctx, cursor, "user:*", 10).Result()
		if err != nil {
			return err
		}

		if newCursor == 0 {
			break
		}
		cursor = newCursor

		if len(keys) != 0 {
			c.db().Del(ctx, keys...)
		}
	}
	//if result.Err() != nil {
	//	return result.Err()
	//}
	//if result.Val() == 0 {
	//	return errors.New("key not found")
	//}
	//result = c.db().HSet(ctx, key, "online", 0)
	//
	//return result.Err()
	return nil
}

func ftSearchConvert[T any](data any) (int64, []T) {
	sliceData, ok := data.([]any)
	if !ok {
		return 0, nil
	}

	count := sliceData[0].(int64)
	result := make([]T, 0, count)

	for i := 1; i < len(sliceData); i += 2 {
		_ = sliceData[i].(string) // Json Key
		values := sliceData[i+1].([]any)
		rawJson := values[len(values)-1].(string)

		var t T
		err := json.Unmarshal([]byte(rawJson), &t)
		if err != nil {
			return 0, nil
		}
		result = append(result, t)
	}
	return count, result
}
