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

func NewUserResponsePayload(self *Client, clients []*Client) Payload {
	payload := Payload{
		Type:     PayloadGetUsers,
		ClientId: self.Id,
	}

	userResponses := make([]UserResponse, 0, len(clients))
	for _, c := range clients {
		// Prevent to response itself
		if c.UserId == self.UserId {
			continue
		}

		userResponses = append(userResponses, UserResponse{
			UserId:   c.UserId,
			Username: c.Username,
			Role:     c.Role,
		})
	}
	payload.Data = userResponses

	return payload
}

type GetUserPayload struct {
	Username string
}
