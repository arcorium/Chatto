package util

import (
	"chatto/internal/constant"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func NewTimeoutContext(parent ...context.Context) (context.Context, context.CancelFunc) {
	if len(parent) == 0 {
		return context.WithTimeout(context.Background(), constant.CONTEXT_TIMEOUT)
	}
	return context.WithTimeout(parent[0], constant.CONTEXT_TIMEOUT)
}

func GetContextValue[T any](key string, ctx *gin.Context) (T, error) {
	data, exist := ctx.Get(key)
	if !exist {
		return *new(T), errors.New("value not found")
	}
	value, ok := data.(T)
	if !ok {
		return *new(T), errors.New("value have different type")
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
