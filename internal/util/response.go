package util

import (
	"server_client_chat/internal/model"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(ctx *gin.Context, code uint, data any) {
	ctx.JSON(int(code), model.NewSuccessResponse(code, data))
}

func ErrorResponse(ctx *gin.Context, code uint, message string) {
	ctx.JSON(int(code), model.NewErrorResponse(code, message, nil))
}
