package model

import "time"

type Role string

const (
	UserRole  Role = "user"
	AdminRole      = "admin"
)

type User struct {
	Id       string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     Role   `json:"role" gorm:"default:user"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

func NewUserResponse(user *User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Name,
		Role:     user.Role,
	}
}
