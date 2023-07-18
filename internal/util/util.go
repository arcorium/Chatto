package util

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const contextTimeout = time.Second * 10

func NewTimeoutContext(parent ...context.Context) (context.Context, context.CancelFunc) {
	if len(parent) == 0 {
		return context.WithTimeout(context.Background(), contextTimeout)
	}
	return context.WithTimeout(parent[0], contextTimeout)
}

func GetContextValue[T any](key string, ctx *gin.Context) (*T, error) {
	var t *T
	data, exist := ctx.Get(key)
	if !exist {
		return t, errors.New("value not found")
	}
	value, ok := data.(*T)
	if !ok {
		return t, errors.New("value have different type")
	}

	return value, nil
}

func HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result), err
}

func ValidatePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
