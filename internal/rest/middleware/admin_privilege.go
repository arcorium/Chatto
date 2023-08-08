package middleware

import (
	"net/http"

	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/util"
	"chatto/internal/util/httputil"
	"github.com/gin-gonic/gin"
)

type AdminPrivilegeMiddleware struct {
}

func (a *AdminPrivilegeMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := util.GetContextValue[model.AccessTokenClaims](constant.KEY_JWT_CLAIMS, c)
		if err != nil || claims.Role != model.AdminRole {
			httputil.ErrorResponse(c, http.StatusForbidden, common.NewError(common.AUTH_UNAUTHORIZED, constant.MSG_AUTH_UNAUTHORIZED))
			c.Abort()
			return
		}
	}
}
