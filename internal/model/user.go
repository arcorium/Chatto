package model

import (
	"chatto/internal/util/strutil"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole      = "admin"
)

func NewUser(name, email, password string, roles ...Role) User {
	var role Role
	if len(roles) > 0 {
		role = roles[0]
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}
	}

	return User{
		Id:             uuid.NewString(),
		Name:           name,
		Email:          email,
		Password:       string(passwordByte),
		Role:           role,
		EmailConfirmed: false,
		CreatedAt:      time.Now(),
	}
}

type User struct {
	Id             string `gorm:"primaryKey;type:uuid"`
	Name           string `gorm:"not null;unique"`
	Email          string `gorm:"not null;unique"`
	Password       string `gorm:"not null"`
	Role           Role   `gorm:"default:user"`
	EmailConfirmed bool   `gorm:"default:false"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (u *User) ValidatePassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword))
}

func (u *User) Validate() bool {
	return !(strutil.IsEmpty(u.Id) || strutil.IsEmpty(u.Name) || strutil.IsEmpty(u.Email) || strutil.IsEmpty(u.Password) || strutil.IsEmpty(string(u.Role)))
}
