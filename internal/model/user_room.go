package model

import (
	"time"
)

func NewUserRoom(roomId string, userId string, role RoomRole) UserRoom {
	return UserRoom{
		RoomId:    roomId,
		UserId:    userId,
		UserRole:  role,
		CreatedAt: time.Now(),
	}
}

type UserRoom struct {
	Id        uint      `gorm:"primaryKey"`
	RoomId    string    `gorm:"not null;type:uuid;uniqueIndex:idx_name"`
	UserId    string    `gorm:"not null;type:uuid;uniqueIndex:idx_name"`
	UserRole  RoomRole  `gorm:"not null;type:text;default:user"`
	CreatedAt time.Time // Used to know when the UserId joining into RoomId
}
