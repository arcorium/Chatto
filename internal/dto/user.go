package dto

import (
	"time"

	"chatto/internal/model"
	"chatto/internal/util/strutil"
)

func NewUserFromCreateInput(input *CreateUserInput) model.User {
	if strutil.IsEmpty(string(input.Role)) {
		input.Role = model.UserRole
	}
	return model.NewUser(input.Username, input.Email, input.Password, input.Role)
}

func NewUserFromUpdateInput(input *UpdateUserInput) model.User {
	user := model.NewUser(input.Username, input.Email, input.Password)
	user.UpdatedAt = time.Now()
	return user
}

type CreateUserInput struct {
	Username string     `json:"username" binding:"required"`
	Email    string     `json:"email" binding:"required"`
	Password string     `json:"password" binding:"required"`
	Role     model.Role `json:"role"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id   string     `json:"id"`
	Name string     `json:"username"`
	Role model.Role `json:"role"`
	// TODO: Add status, so when user search the name it will know either the user is online or no
}

func NewUserResponse(user *model.User) UserResponse {
	return UserResponse{
		Id:   user.Id,
		Name: user.Name,
		Role: user.Role,
	}
}

type GetUserInput struct {
	Username string `json:"username"`
}

func NewGetUserOutput(userResponses []UserResponse) GetUserOutput {
	return GetUserOutput{
		Users: userResponses,
	}
}

type GetUserOutput struct {
	Users []UserResponse `json:"users"`
}
