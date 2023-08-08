package model

import (
	"github.com/google/uuid"
	"time"
)

type TokenType string

func NewCredential(userId string, deviceName string, token string) Credential {
	return Credential{
		Id:         uuid.NewString(),
		UserId:     userId,
		DeviceName: deviceName,
		Token:      token,
	}
}

type Credential struct {
	Id         string `json:"id" gorm:"primaryKey;type:uuid"`
	UserId     string `json:"user_id" gorm:"not null;type:uuid"`
	DeviceName string `json:"device_name"`
	Token      string `json:"token"`

	CreatedAt time.Time
}

type AccessTokenClaims struct {
	UserId    string
	Name      string
	Role      Role
	RefreshId string
	//Exp       uint64
}

type Device struct {
	Id         string `json:"id"`
	DeviceName string `json:"device"`
}
