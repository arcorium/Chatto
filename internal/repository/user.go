package repository

import (
	"chatto/internal/model"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) FindUsers() ([]model.User, error) {
	var users []model.User
	result := u.db.Find(&users)
	return users, result.Error
}

func (u *UserRepository) FindUserById(id string) (*model.User, error) {
	var user model.User
	result := u.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (u *UserRepository) FindUserByName(name string) (*model.User, error) {
	var user model.User
	result := u.db.First(&user, "name = ?", name)
	return &user, result.Error
}

func (u *UserRepository) UpdateUserById(id string, user *model.User) error {
	result := u.db.Model(&model.User{Id: id}).Updates(user)
	return result.Error
}

func (u *UserRepository) CreateUser(user *model.User) error {
	result := u.db.Create(user)
	return result.Error
}

func (u *UserRepository) RemoveUserById(id string) error {
	result := u.db.Delete(&model.User{}, "id = ?", id)
	return result.Error
}
