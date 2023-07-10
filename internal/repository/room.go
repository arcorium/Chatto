package repository

import (
	"chatto/internal/model"
	"gorm.io/gorm"
)

func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &roomRepository{db: db}
}

type roomRepository struct {
	db *gorm.DB
}

func (r roomRepository) FindRooms() ([]model.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r roomRepository) FindRoomById(id string) (model.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r roomRepository) FindUsersInRoomById(id string) ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r roomRepository) RemoveRoomById(id string) (model.Room, error) {
	//TODO implement me
	panic("implement me")
}
