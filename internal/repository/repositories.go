package repository

import (
	"chatto/internal/dto"
	"chatto/internal/model"
)

type IUserRepository interface {
	FindUsers() ([]model.User, error)
	FindUserById(id string) (*model.User, error)
	FindUserByName(name string) (*model.User, error)
	FindUsersByLikelyName(name string) ([]model.User, error)
	UpdateUserById(id string, user *model.User) error
	CreateUser(user *model.User) error
	RemoveUserById(id string) error
}

type IAuthRepository interface {
	FindTokenById(tokenId string) (model.Credential, error)
	FindTokenByUserId(userId string) ([]model.Credential, error)
	FindDevicesByUserId(userId string) ([]model.Device, error)
	CreateToken(token *model.Credential) error
	UpdateToken(originalId string, token *model.Credential) error
	RemoveTokenById(tokenId string) error
	RemoveTokensByUserId(userId string) error
	RemoveAllToken() error
}

type IRoomRepository interface {
	CreateRoom(room *model.Room) error
	FindRooms() ([]model.Room, error)
	FindRoomById(roomId string) (*model.Room, error)
	DeleteRoomById(roomId string) error
	FindRoomsByUserId(userId string) ([]model.Room, error)
}

type IUserRoomRepository interface {
	GetUserIdsOnRoomById(roomId string) ([]string, error)
	GetRoomMemberCountById(roomId string) (int64, error)
	FindUserRoomsByUserId(userId string) ([]model.UserRoom, error)
	FindUsersByRoomId(roomId string) ([]model.User, error)
	AddUsersIntoRoomById(userRoom []model.UserRoom) error
	RemoveUserFromRoomById(roomId string, userId string) error
	RemoveAllUsersFromRoomById(roomId string) error
	RemoveUsersFromRoomById(roomId string, userId []string) error
	RemoveAllRoomsFromUserById(userId string) error
}

type IChatRepository interface {
	// CreateMessage Will store new message
	CreateMessage(message *model.Message) error
	// CreateNotification Will store new notification
	CreateNotification(notif *model.Notification) error
	// FindRoomChats Used to get all chats based on the roomId and range time
	FindRoomChats(request *dto.MessageRequest) ([]model.Message, error)
	// FindRoomNotifications Used to get all notifications based on the roomId and range time
	FindRoomNotifications(request *dto.NotificationRequest) ([]model.Notification, error)
	// NewClient Used to either create new key or increment "online" key by 1
	NewClient(client *model.Client) error
	// RemoveClient Decrement "online" key by 1
	RemoveClient(client *model.Client) error
	// ResetClients Used to reset online field to 0
	ResetClients() error
}
