package repository

import (
	"chatto/internal/model"

	"gorm.io/gorm"
)

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db: db}
}

type authRepository struct {
	db *gorm.DB
}

func (a *authRepository) FindTokenById(tokenId string) (model.TokenDetails, error) {
	var token model.TokenDetails
	result := a.db.First(&token, "id = ?", tokenId)

	return token, result.Error
}

func (a *authRepository) FindTokenByUserId(userId string) ([]model.TokenDetails, error) {
	var token []model.TokenDetails
	result := a.db.Find(&token, "user_id = ?", userId)

	return token, result.Error
}

func (a *authRepository) SaveToken(token *model.TokenDetails) error {
	result := a.db.Create(token) // Not updating record, just create
	return result.Error
}

func (a *authRepository) RemoveTokenById(tokenId string) error {
	result := a.db.Delete(&model.TokenDetails{}, "id = ?", tokenId)
	return result.Error
}

func (a *authRepository) RemoveTokensByUserId(userId string) error {
	result := a.db.Delete(&model.TokenDetails{}, "user_id = ?", userId)
	return result.Error
}
