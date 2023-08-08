package service

import (
	"net/http"

	"chatto/internal/dto"
	"chatto/internal/repository"
	"chatto/internal/util/containers"

	"chatto/internal/constant"
	"chatto/internal/model/common"
)

type IUserService interface {
	GetUsers() ([]dto.UserResponse, common.Error)
	FindUserById(id string) (dto.UserResponse, common.Error)
	FindUsersByLikelyName(name string) ([]dto.UserResponse, common.Error)
	FindAndValidateUserByName(name, password string) (dto.UserResponse, common.Error)
	UpdateUserById(id string, user *dto.UpdateUserInput) common.Error
	CreateUser(user *dto.CreateUserInput) common.Error
	RemoveUserById(id string) common.Error
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

type userService struct {
	userRepo repository.IUserRepository
}

func (u userService) FindUsersByLikelyName(name string) ([]dto.UserResponse, common.Error) {
	users, err := u.userRepo.FindUsersByLikelyName(name)
	if err != nil {
		return nil, common.NewError(common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
	}
	userResponses := containers.ConvertSlice(users, dto.NewUserResponse)

	return userResponses, common.NoError()
}

func (u userService) FindUserById(id string) (dto.UserResponse, common.Error) {
	user, err := u.userRepo.FindUserById(id)
	return dto.NewUserResponse(user), common.NewConditionalError(err, common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
}

func (u userService) FindAndValidateUserByName(name, password string) (dto.UserResponse, common.Error) {
	user, err := u.userRepo.FindUserByName(name)
	if err != nil {
		return dto.UserResponse{}, common.NewError(common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
	}

	err = user.ValidatePassword(password)
	return dto.NewUserResponse(user), common.NewConditionalError(err, common.AUTH_PASSWORD_INVALID_ERROR, constant.MSG_FAILED_USER_LOGIN)
}

func (u userService) UpdateUserById(id string, input *dto.UpdateUserInput) common.Error {
	user := dto.NewUserFromUpdateInput(input)
	err := u.userRepo.UpdateUserById(id, &user)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_FAILED_UPDATE_USER)
}

func (u userService) GetUsers() ([]dto.UserResponse, common.Error) {
	users, err := u.userRepo.FindUsers()
	if err != nil {
		return nil, common.NewError(http.StatusBadRequest, constant.MSG_USER_NOT_FOUND)
	}

	userResponses := containers.ConvertSlice(users, dto.NewUserResponse)

	return userResponses, common.NoError()
}

func (u userService) CreateUser(input *dto.CreateUserInput) common.Error {
	user := dto.NewUserFromCreateInput(input)
	if !user.Validate() {
		return common.NewError(common.USER_CREATION_ERROR, constant.MSG_CREATE_USER_FAILED)
	}
	err := u.userRepo.CreateUser(&user)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_CREATE_USER_FAILED)
}

func (u userService) RemoveUserById(id string) common.Error {
	err := u.userRepo.RemoveUserById(id)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_FAILED_REMOVE_USER)
}
