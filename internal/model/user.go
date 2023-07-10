package model

import "time"

type Role string

const (
	UserRole  Role = "user"
	AdminRole      = "admin"
)

type User struct {
	Id       string   `json:"id" gorm:"primaryKey;type:uuid"`
	Name     string   `json:"name" gorm:"not null"`
	Password string   `json:"password" gorm:"not null"`
	Role     Role     `json:"role" gorm:"default:user"`
	RoomIds  []string `json:"room_ids" gorm:"type:uuid[]"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type UserResponse struct {
	UserId   string `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

func NewUserResponse(user *User) UserResponse {
	return UserResponse{
		UserId:   user.Id,
		Username: user.Name,
		Role:     user.Role,
	}
}

type IncomeGetUser struct {
	Username string `json:"username"`
}

func NewOutcomeGetUser(clients []*Client) OutcomeGetUser {
	userResponses := make([]UserResponse, 0, len(clients))
	for _, c := range clients {
		userResponses = append(userResponses, UserResponse{
			UserId:   c.UserId,
			Username: c.Username,
			Role:     c.Role,
		})
	}
	return OutcomeGetUser{
		Users: userResponses,
	}
}

type OutcomeGetUser struct {
	Users []UserResponse `json:"users"`
}
