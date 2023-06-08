package service

import (
	"errors"

	"chatto/internal/model"
)

type CustomError struct {
	HttpCode  uint
	ErrorCode uint
	Err       error
}

func (e CustomError) IsError() bool {
	return e.Err != nil
}

func (e CustomError) Error() string {
	return e.Err.Error()
}

func NewError(code uint, message string) CustomError {
	return CustomError{
		HttpCode:  code,
		ErrorCode: code,
		Err:       errors.New(message),
	}
}

func NoError() CustomError {
	return CustomError{
		HttpCode:  200,
		ErrorCode: 200,
		Err:       nil,
	}
}

type IUserService interface {
	FindUsers() ([]model.UserResponse, CustomError)
	FindUserById(id string) (model.UserResponse, CustomError)
	UpdateUserById(id string, user *model.User) CustomError
	CreateUser(user *model.User) CustomError
	RemoveUserById(id string) CustomError
}

type IAuthService interface {
	SignIn(input *model.SignInInput, sysInfo *model.SystemInfo) (string, CustomError)
	SignUp(input *model.SignInInput) CustomError
	Logout(userId string, tokenId string) CustomError
	LogoutAllDevice(userId string) CustomError
	RefreshToken(accessToken string) (string, CustomError)
}

type IChatService interface {
	HandleNewClient(client *model.Client) CustomError
}
