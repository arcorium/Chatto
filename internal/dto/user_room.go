package dto

import (
	"chatto/internal/model"
)

func NewUserRole(userId string, role model.RoomRole) UserWithRole {
	return UserWithRole{
		UserId: userId,
		Role:   role,
	}
}

// NewUserRoles Used to create slice of UserWithRole with set all the user role as model.RoomRole
func NewUserRoles(role model.RoomRole, userIds ...string) []UserWithRole {
	users := make([]UserWithRole, 0, len(userIds))
	for _, userId := range userIds {
		users = append(users, NewUserRole(userId, role))
	}
	return users
}

type UserWithRole struct {
	UserId string         `json:"user_id"`
	Role   model.RoomRole `json:"room_role"`
}

type UserRoomAddInput struct {
	Users  []UserWithRole `json:"users"`
	RoomId string         `json:"room_id"`
}
type UserRoomRemoveInput struct {
	UserIds []string `json:"user_ids"`
	RoomId  string   `json:"room_id"`
}

func NewUserRoomResponse(room *model.UserRoom) UserRoomResponse {
	return UserRoomResponse{
		RoomId:   room.RoomId,
		UserRole: room.UserRole,
	}
}

type UserRoomResponse struct {
	RoomId   string         `json:"room_id"`
	UserRole model.RoomRole `json:"user_role"`
}
