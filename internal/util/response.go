package util

import (
	"chatto/internal/model/common"
	"github.com/gin-gonic/gin"
)

func SuccessResponse(ctx *gin.Context, code uint, data any) {
	ctx.JSON(int(code), common.NewSuccessResponse(code, data))
}

func ErrorResponse(ctx *gin.Context, code uint, message string) {
	ctx.JSON(int(code), common.NewErrorResponse(code, message, nil))
}
