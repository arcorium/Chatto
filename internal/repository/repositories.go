package repository

import (
	"chatto/internal/model"
)

type IUserRepository interface {
	FindUsers() ([]model.User, error)
	FindUserById(id string) (*model.User, error)
	FindUserByName(name string) (*model.User, error)
	UpdateUserById(id string, user *model.User) error
	CreateUser(user *model.User) error
	RemoveUserById(id string) error
}

type IAuthRepository interface {
	FindTokenById(id string) (model.TokenDetails, error)
	FindTokenByUserId(id string) ([]model.TokenDetails, error)
	SaveToken(token *model.TokenDetails) error
	RemoveTokenById(id string) error
	RemoveTokensByUserId(userId string) error
}

type IRoomRepository interface {
	FindRooms() ([]model.Room, error)
	FindRoomById(id string) (model.Room, error)
	FindUsersInRoomById(id string) ([]model.User, error)
	RemoveRoomById(id string) (model.Room, error)
}

type IClientRepository interface {
	RegisterClient()
}

type IChatRepository interface {
	InsertChat(message *model.Message) error
	FindChats() ([]model.Message, error)
}
