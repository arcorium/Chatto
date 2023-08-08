package pg_repo

import (
	"chatto/internal/model"
	"chatto/internal/repository"
	"gorm.io/gorm"
)

func NewRoomRepository(db *gorm.DB) repository.IRoomRepository {
	return &roomRepository{db_: db}
}

type roomRepository struct {
	db_ *gorm.DB
}

func (r roomRepository) db() *gorm.DB {
	return r.db_.Debug()
}

func (r roomRepository) CreateRoom(room *model.Room) error {
	result := r.db().Create(room)
	return result.Error
}

func (r roomRepository) FindRooms() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db().Find(&rooms)
	return rooms, result.Error
}

func (r roomRepository) FindRoomById(roomId string) (*model.Room, error) {
	var room model.Room
	result := r.db().Find(&room, "id = ?", roomId)
	return &room, result.Error
}

func (r roomRepository) FindRoomsByUserId(userId string) ([]model.Room, error) {
	var rooms []model.Room
	result := r.db().Raw("SELECT rooms.* FROM user_rooms INNER JOIN rooms ON rooms.id = user_rooms.room_id WHERE user_rooms.user_id = ?", userId).Scan(&rooms)
	return rooms, result.Error
}

func (r roomRepository) DeleteRoomById(roomId string) error {
	result := r.db().Delete(&model.Room{}, "id = ?", roomId)
	return result.Error
}
