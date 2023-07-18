package service

import (
	"chatto/internal/dto"
	"chatto/internal/model"
	"chatto/internal/repository"
	"chatto/internal/util/ctrutil"
	"net/http"

	"chatto/internal/constant"
	"chatto/internal/model/common"

	"chatto/internal/util"

	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers() ([]dto.UserResponse, common.Error)
	FindUserById(id string) (dto.UserResponse, common.Error)
	GetUsersOnRoomById(roomId string) ([]dto.UserResponse, common.Error)
	UpdateUserById(id string, user *dto.UpdateUserInput) common.Error
	CreateUser(user *dto.CreateUserInput) common.Error
	RemoveUserById(id string) common.Error
}

func NewUserService(userRepo repository.IUserRepository, userRoomRepo repository.IUserRoomRepository) IUserService {
	return &userService{userRepo: userRepo, userRoomRepo: userRoomRepo}
}

type userService struct {
	userRepo     repository.IUserRepository
	userRoomRepo repository.IUserRoomRepository
}

func (u userService) FindUserById(id string) (dto.UserResponse, common.Error) {
	user, err := u.userRepo.FindUserById(id)
	return dto.NewUserResponse(user), common.NewConditionalError(err, common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
}

func (u userService) UpdateUserById(id string, input *dto.UpdateUserInput) common.Error {
	user := dto.NewUserFromUpdateInput(input)
	err := u.userRepo.UpdateUserById(id, &user)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_FAILED_UPDATE_USER)
}

func (u userService) GetUsers() ([]dto.UserResponse, common.Error) {
	users, err := u.userRepo.FindUsers()
	if err != nil {
		return nil, common.NewError(http.StatusBadRequest, constant.MSG_USER_NOT_FOUND)
	}

	userResponses := ctrutil.ConvertSliceType(users, func(current *model.User) dto.UserResponse {
		return dto.NewUserResponse(current)
	})

	return userResponses, common.NoError()
}

func (u userService) GetUsersOnRoomById(roomId string) ([]dto.UserResponse, common.Error) {
	userIds, err := u.userRoomRepo.GetUserIdsOnRoomById(roomId)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	userResponses, err := ctrutil.SafeConvertSliceType(userIds, func(current *string) (dto.UserResponse, error) {
		user, err := u.userRepo.FindUserById(*current)
		return dto.NewUserResponse(user), err
	})
	return userResponses, common.NewConditionalError(err, common.USER_NOT_FOUND_ERROR, constant.MSG_USER_NOT_FOUND)
}

func (u userService) CreateUser(input *dto.CreateUserInput) common.Error {
	user := dto.NewUserFromCreateInput(input)
	user.Id = uuid.NewString()
	password, err := util.HashPassword(input.Password)
	if err != nil {
		return common.NewError(common.HASH_PASSWORD_ERROR, constant.MSG_FAILED_CREATE_USER)
	}
	user.Password = password

	err = u.userRepo.CreateUser(&user)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_FAILED_CREATE_USER)
}

func (u userService) RemoveUserById(id string) common.Error {
	err := u.userRepo.RemoveUserById(id)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_FAILED_REMOVE_USER)
}
