package pg_repo

import (
	"chatto/internal/model"
	"chatto/internal/repository"
	"gorm.io/gorm"
)

func NewUserRoomRepository(db *gorm.DB) repository.IUserRoomRepository {
	return &userRoomRepository{db_: db}
}

type userRoomRepository struct {
	db_ *gorm.DB
}

func (u userRoomRepository) db() *gorm.DB {
	return u.db_.Debug()
}

func (u userRoomRepository) GetUserIdsOnRoomById(roomId string) ([]string, error) {
	var userIds []string
	result := u.db().Find(&userIds, "room_id = ?", roomId)
	return userIds, result.Error
}

func (u userRoomRepository) GetRoomIdsByUserId(userId string) ([]string, error) {
	var roomIds []string
	result := u.db().Find(&roomIds, "user_id = ?", userId)
	return roomIds, result.Error
}

func (u userRoomRepository) AddUserIntoRoomById(userRoom *model.UserRoom) error {
	result := u.db().Create(*userRoom)
	return result.Error
}

func (u userRoomRepository) AddUsersIntoRoomById(userRoom []model.UserRoom) error {
	return u.db().Transaction(func(tx *gorm.DB) error {
		for _, data := range userRoom {
			result := tx.Create(data)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (u userRoomRepository) RemoveUserFromRoomById(roomId string, userId string) error {
	result := u.db().Where("room_id = ? AND user_id = ?", roomId, userId).Delete(&model.UserRoom{})
	return result.Error
}

func (u userRoomRepository) RemoveAllUsersFromRoomById(roomId string) error {
	result := u.db().Where("room_id = ?", roomId).Delete(&model.UserRoom{})
	return result.Error
}

func (u userRoomRepository) RemoveUsersFromRoomById(roomId string, userId []string) error {
	mdl := &model.UserRoom{}
	return u.db().Transaction(func(tx *gorm.DB) error {
		for _, id := range userId {
			result := tx.Where("room_id = ? AND user_id = ?", roomId, id).Delete(mdl)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (u userRoomRepository) RemoveAllRoomsFromUserById(userId string) error {
	result := u.db().Where("user_id = ?", userId).Delete(&model.UserRoom{})
	return result.Error
}
