package model

import (
	"time"
)

func NewUserRoom(roomId string, userId string) UserRoom {
	return UserRoom{
		RoomId: roomId,
		UserId: userId,
	}
}

type UserRoom struct {
	Id        uint      `gorm:"primarykey"`
	RoomId    string    `gorm:"not null"`
	UserId    string    `gorm:"not null"`
	CreatedAt time.Time // Used to know when the UserId joining into RoomId
}
