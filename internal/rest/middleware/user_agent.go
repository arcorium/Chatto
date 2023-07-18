package middleware

import (
	"chatto/internal/constant"
	"chatto/internal/model/common"
	"chatto/internal/util/httputil"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

type UserAgentValidationConfig struct {
}
type UserAgentValidationMiddleware struct {
	Config *UserAgentValidationConfig
}

func (u *UserAgentValidationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if len(userAgent) == 0 {
			httputil.ErrorResponse(c, http.StatusBadRequest, common.NewError(common.USER_AGENT_UNKNOWN_ERROR, constant.MSG_USER_AGENT_UNKNOWN))
			c.Abort()
			return
		}

		systemInfo := u.parseUserAgent(userAgent)

		c.Set(constant.KEY_USER_AGENT, &systemInfo)
		c.Next()
	}
}

func (u *UserAgentValidationMiddleware) parseUserAgent(userAgent string) common.SystemInfo {
	info := useragent.Parse(userAgent)
	return common.SystemInfo{
		Name: info.Name,
		Os:   info.OS,
	}
}
