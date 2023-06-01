package repository

import (
	"server_client_chat/internal/model"

	"gorm.io/gorm"
)

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return AuthRepository{db: db}
}

type AuthRepository struct {
	db *gorm.DB
}

func (a *AuthRepository) FindTokenById(tokenId string) (model.TokenDetails, error) {
	var token model.TokenDetails
	result := a.db.First(&token, "id = ?", tokenId)

	return token, result.Error
}

func (a *AuthRepository) FindTokenByUserId(userId string) ([]model.TokenDetails, error) {
	var token []model.TokenDetails
	result := a.db.Find(&token, "user_id = ?", userId)

	return token, result.Error
}

func (a *AuthRepository) SaveToken(token *model.TokenDetails) error {
	result := a.db.Create(token) // Not updating record, just create
	return result.Error
}

func (a *AuthRepository) RemoveTokenById(tokenId string) error {
	result := a.db.Delete(&model.TokenDetails{}, "id = ?", tokenId)
	return result.Error
}

func (a *AuthRepository) RemoveTokensByUserId(userId string) error {
	result := a.db.Delete(&model.TokenDetails{}, "user_id = ?", userId)
	return result.Error
}
