package service

import (
	"log"
	"net/http"

	"server_client_chat/internal/config"
	"server_client_chat/internal/model"
	"server_client_chat/internal/repository"
	"server_client_chat/internal/util"

	"github.com/google/uuid"
)

func NewUserService(serverConfig *config.AppConfig, repository repository.IUserRepository) UserService {
	return UserService{ServerConfig: serverConfig, repo: repository}
}

type UserService struct {
	ServerConfig *config.AppConfig
	repo         repository.IUserRepository
}

func (u *UserService) FindUserById(id string) (model.UserResponse, CustomError) {
	user, err := u.repo.FindUserById(id)
	if err != nil {
		log.Println(err)
		return model.UserResponse{}, NewError(http.StatusBadRequest, util.ERR_USER_NOT_FOUND)
	}
	return model.NewUserResponse(user), NoError()
}

func (u *UserService) UpdateUserById(id string, user *model.User) CustomError {
	err := u.repo.UpdateUserById(id, user)
	if err != nil {
		log.Println(err)
		return NewError(http.StatusBadRequest, util.ERR_USER_UPDATE)
	}
	return NoError()
}

func (u *UserService) FindUsers() ([]model.UserResponse, CustomError) {
	users, err := u.repo.FindUsers()
	if err != nil {
		log.Println(err)
		return nil, NewError(http.StatusBadRequest, util.ERR_USER_NOT_FOUND)
	}

	usersResponse := make([]model.UserResponse, 0)
	for _, x := range users {
		usersResponse = append(usersResponse, model.NewUserResponse(&x))
	}

	return usersResponse, NoError()
}

func (u *UserService) CreateUser(user *model.User) CustomError {
	user.Id = uuid.NewString()
	password, err := util.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		return NewError(http.StatusBadRequest, util.ERR_USER_CREATE)
	}
	user.Password = password

	err = u.repo.CreateUser(user)
	return NoError()
}

func (u *UserService) RemoveUserById(id string) CustomError {
	err := u.repo.RemoveUserById(id)
	if err != nil {
		return NewError(http.StatusBadRequest, util.ERR_USER_REMOVE)
	}
	return NoError()
}
