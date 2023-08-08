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

func (u userRoomRepository) GetRoomMemberCountById(roomId string) (int64, error) {
	var count int64
	res := u.db().Model(&model.UserRoom{}).Where("room_id = ?", roomId).Count(&count)
	return count, res.Error
}

func (u userRoomRepository) FindUserRoomsByUserId(userId string) ([]model.UserRoom, error) {
	// TODO: Join with room table
	var userRooms []model.UserRoom
	err := u.db().Find(&userRooms, "user_id = ?", userId)
	return userRooms, err.Error
}

func (u userRoomRepository) FindUsersByRoomId(roomId string) ([]model.User, error) {
	// TODO: Join with user table
	var users []model.User
	result := u.db().Raw("SELECT users.* from user_rooms INNER JOIN users ON user_rooms.user_id = users.id WHERE user_rooms.room_id = ?", roomId).Scan(&users)
	return users, result.Error
}

func (u userRoomRepository) GetUserIdsOnRoomById(roomId string) ([]string, error) {
	var userIds []string
	result := u.db().Find(&userIds, "room_id = ?", roomId)
	return userIds, result.Error
}

func (u userRoomRepository) AddUsersIntoRoomById(userRoom []model.UserRoom) error {
	return u.db().Transaction(func(tx *gorm.DB) error {
		for _, data := range userRoom {
			result := tx.Create(&data)
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
