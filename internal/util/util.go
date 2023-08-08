package util

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func EscapeMinesSymbols(str string) string {
	return strings.ReplaceAll(str, "-", "\\-")
}
