package dto

import (
	"chatto/internal/model"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewSignInOutput(jwtToken string) SignInOutput {
	return SignInOutput{
		Type:  "Bearer",
		Token: jwtToken,
	}
}

type SignInOutput struct {
	Type  string `json:"type"`
	Token string `json:"access_token"`
}

func NewUserFromSignUpInput(input *SignUpInput) model.User {
	return model.NewUser(input.Username, input.Email, input.Password, model.UserRole)
}

type SignUpInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenInput struct {
	Type        string `json:"type" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

func NewRefreshTokenOutput(jwtToken string) RefreshTokenOutput {
	return RefreshTokenOutput{
		Type:  "Bearer",
		Token: jwtToken,
	}
}

type RefreshTokenOutput struct {
	Type  string `json:"type"`
	Token string `json:"access_token"`
}
