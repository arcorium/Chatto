package repository

import (
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
}

type IUserRoomRepository interface {
	GetUserIdsOnRoomById(roomId string) ([]string, error)
	GetRoomIdsByUserId(userId string) ([]string, error)
	AddUserIntoRoomById(userRoom *model.UserRoom) error
	AddUsersIntoRoomById(userRoom []model.UserRoom) error
	RemoveUserFromRoomById(roomId string, userId string) error
	RemoveAllUsersFromRoomById(roomId string) error
	RemoveUsersFromRoomById(roomId string, userId []string) error
	RemoveAllRoomsFromUserById(userId string) error
}

type IChatRepository interface {
	CreateMessage(id string, message *model.Message) error
	CreateNotification(id string, notif *model.Notification) error
	FindChats() ([]model.Message, error)
	NewClient(client *model.Client) error
	RemoveClient(client *model.Client) error
}
