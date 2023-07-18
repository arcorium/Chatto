package pg_repo

import (
	"chatto/internal/model"
	"chatto/internal/repository"

	"gorm.io/gorm"
)

func NewAuthRepository(db *gorm.DB) repository.IAuthRepository {
	return &authRepository{db_: db}
}

type authRepository struct {
	db_ *gorm.DB
}

func (a *authRepository) db() *gorm.DB {
	return a.db_.Debug()
}

func (a *authRepository) FindTokenById(tokenId string) (model.Credential, error) {
	var token model.Credential
	result := a.db().First(&token, "id = ?", tokenId)

	return token, result.Error
}

func (a *authRepository) FindTokenByUserId(userId string) ([]model.Credential, error) {
	var token []model.Credential
	result := a.db().Find(&token, "user_id = ?", userId)

	return token, result.Error
}
func (a *authRepository) FindDevicesByUserId(id string) ([]model.Device, error) {
	var devices []model.Device
	result := a.db().Model(&model.Credential{}).Select("device_name, id").Where("user_id = ?", id).Scan(&devices)
	return devices, result.Error
}

func (a *authRepository) CreateToken(token *model.Credential) error {
	result := a.db().Create(token) // Not updating record, just create
	return result.Error
}

func (a *authRepository) UpdateToken(originalId string, token *model.Credential) error {
	result := a.db().Where("id = ?", originalId).Updates(token)
	return result.Error
}

func (a *authRepository) RemoveTokenById(tokenId string) error {
	result := a.db().Delete(&model.Credential{}, "id = ?", tokenId)
	return result.Error
}

func (a *authRepository) RemoveTokensByUserId(userId string) error {
	result := a.db().Delete(&model.Credential{}, "user_id = ?", userId)
	return result.Error
}

func (a *authRepository) RemoveAllToken() error {
	result := a.db().Delete(&model.Credential{})
	return result.Error
}
