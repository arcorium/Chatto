package service

import (
	"chatto/internal/model"
	"chatto/internal/model/common"
)

type IUserService interface {
	FindUsers() ([]model.UserResponse, common.Error)
	FindUserById(id string) (model.UserResponse, common.Error)
	UpdateUserById(id string, user *model.User) common.Error
	CreateUser(user *model.User) common.Error
	RemoveUserById(id string) common.Error
}

type IAuthService interface {
	SignIn(input *model.SignInInput, sysInfo *common.SystemInfo) (string, common.Error)
	SignUp(input *model.SignInInput) common.Error
	Logout(userId string, tokenId string) common.Error
	LogoutAllDevice(userId string) common.Error
	RefreshToken(accessToken string) (string, common.Error)
}

type IRoomService interface {
	CreateRoom()
	FindRooms()
	FindRoomByName()
	FindRoomById()
	DeleteRoom()
	AddUserIntoRoom()
	DeleteUserFromRoom()
}
