package middleware

import (
	"chatto/internal/constant"
	"chatto/internal/model/common"

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
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(http.StatusBadRequest, constant.ERR_USER_AGENT_MIDDLEWARE, nil))
			return
		}

		systemInfo := u.parseUserAgent(userAgent)

		c.Set(KEY_USER_AGENT, systemInfo)
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
