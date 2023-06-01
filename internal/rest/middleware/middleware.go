package middleware

import "github.com/gin-gonic/gin"

type IMiddleware interface {
	Handle() gin.HandlerFunc
}

const (
	KEY_USER_AGENT = "system"
	KEY_JWT_CLAIMS = "claims"
)
