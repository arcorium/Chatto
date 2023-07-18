package httputil

import (
	"chatto/internal/model/common"
	"github.com/gin-gonic/gin"
)

// SuccessResponse Used to set common.Response as response for success response
func SuccessResponse(ctx *gin.Context, httpCode int, err common.Error, data any) {
	ctx.JSON(httpCode, common.NewSuccessResponse(err, data))
}

// ErrorResponse Used to set common.Response as response for bad response
func ErrorResponse(ctx *gin.Context, httpCode int, err common.Error) {
	ctx.JSON(httpCode, common.NewErrorResponse(err))
}

// ConditionalResponse Used to set response based on the common.Error passed as parameter. depending
// on the value of the common.Error it will call ErrorResponse or SuccessResponse
func ConditionalResponse(ctx *gin.Context, err common.Error, badHttpCode int, successHttpCode int, successData any) {
	if err.IsError() {
		ErrorResponse(ctx, badHttpCode, err)
	} else {
		SuccessResponse(ctx, successHttpCode, err, successData)
	}
}
