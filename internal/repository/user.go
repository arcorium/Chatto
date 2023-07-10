package repository

import (
	"chatto/internal/model"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) FindUsers() ([]model.User, error) {
	var users []model.User
	result := u.db.Find(&users)
	return users, result.Error
}

func (u userRepository) FindUserById(id string) (*model.User, error) {
	var user model.User
	result := u.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (u userRepository) FindUserByName(name string) (*model.User, error) {
	var user model.User
	result := u.db.First(&user, "name = ?", name)
	return &user, result.Error
}

func (u userRepository) UpdateUserById(id string, user *model.User) error {
	result := u.db.Model(&model.User{Id: id}).Updates(user)
	return result.Error
}

func (u userRepository) CreateUser(user *model.User) error {
	result := u.db.Create(user)
	return result.Error
}

func (u userRepository) RemoveUserById(id string) error {
	result := u.db.Delete(&model.User{}, "id = ?", id)
	return result.Error
}
