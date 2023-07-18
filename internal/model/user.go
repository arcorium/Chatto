package model

import "time"

type Role string

const (
	UserRole  Role = "user"
	AdminRole      = "admin"
)

type User struct {
	Id             string `gorm:"primaryKey;type:uuid"`
	Name           string `gorm:"not null"`
	Email          string `gorm:"not null"`
	Password       string `gorm:"not null"`
	Role           Role   `gorm:"default:user"`
	EmailConfirmed bool   `gorm:"default:false"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}
