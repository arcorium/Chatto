package service

import (
	"log"
	"net/http"

	"chatto/internal/constant"
	"chatto/internal/model/common"

	"chatto/internal/model"
	"chatto/internal/repository"
	"chatto/internal/util"

	"github.com/google/uuid"
)

func NewUserService(repository repository.IUserRepository) IUserService {
	return &userService{repo: repository}
}

type userService struct {
	repo repository.IUserRepository
}

func (u userService) FindUserById(id string) (model.UserResponse, common.Error) {
	user, err := u.repo.FindUserById(id)
	if err != nil {
		log.Println(err)
		return model.UserResponse{}, common.NewError(http.StatusBadRequest, constant.ERR_USER_NOT_FOUND)
	}
	return model.NewUserResponse(user), common.NoError()
}

func (u userService) UpdateUserById(id string, user *model.User) common.Error {
	err := u.repo.UpdateUserById(id, user)
	if err != nil {
		log.Println(err)
		return common.NewError(http.StatusBadRequest, constant.ERR_USER_UPDATE)
	}
	return common.NoError()
}

func (u userService) FindUsers() ([]model.UserResponse, common.Error) {
	users, err := u.repo.FindUsers()
	if err != nil {
		log.Println(err)
		return nil, common.NewError(http.StatusBadRequest, constant.ERR_USER_NOT_FOUND)
	}

	usersResponse := make([]model.UserResponse, 0)
	for _, x := range users {
		usersResponse = append(usersResponse, model.NewUserResponse(&x))
	}

	return usersResponse, common.NoError()
}

func (u userService) CreateUser(user *model.User) common.Error {
	user.Id = uuid.NewString()
	password, err := util.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		return common.NewError(http.StatusBadRequest, constant.ERR_USER_CREATE)
	}
	user.Password = password

	err = u.repo.CreateUser(user)
	return common.NoError()
}

func (u userService) RemoveUserById(id string) common.Error {
	err := u.repo.RemoveUserById(id)
	if err != nil {
		return common.NewError(http.StatusBadRequest, constant.ERR_USER_REMOVE)
	}
	return common.NoError()
}
