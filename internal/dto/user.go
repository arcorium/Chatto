package dto

import (
	"chatto/internal/model"
)

func NewUserFromCreateInput(input *CreateUserInput) model.User {
	return model.User{
		Name:     input.Username,
		Email:    input.Email,
		Password: input.Password,
	}
}

func NewUserFromUpdateInput(input *UpdateUserInput) model.User {
	return model.User{
		Name:     input.Username,
		Email:    input.Email,
		Password: input.Password,
	}
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       string     `json:"id"`
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
	// TODO: Add status, so when user search the name it will know either the user is current online or no
}

func NewUserResponse(user *model.User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Name,
		Role:     user.Role,
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
