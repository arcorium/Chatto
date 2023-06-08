package middleware

import (
	"chatto/internal/constant"
	"net/http"

	"chatto/internal/model"
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
			c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, constant.ERR_USER_AGENT_MIDDLEWARE, nil))
			return
		}

		systemInfo := u.parseUserAgent(userAgent)

		c.Set(KEY_USER_AGENT, systemInfo)
		c.Next()
	}
}

func (u *UserAgentValidationMiddleware) parseUserAgent(userAgent string) model.SystemInfo {
	info := useragent.Parse(userAgent)
	return model.SystemInfo{
		Name: info.Name,
		Os:   info.OS,
	}
}
