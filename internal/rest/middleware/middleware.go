package middleware

import (
	"chatto/internal/config"
	"github.com/gin-gonic/gin"
)

func NewMiddleware(config *config.AppConfig) Middleware {
	tokenValConf := TokenValidationConfig{
		SecretKeyFunc: config.JWTKeyFunc,
		TokenType:     "Bearer",
		SigningType:   config.JWTSigningType,
	}
	userAgentValConf := UserAgentValidationConfig{}

	middleware := TokenValidationMiddleware{Config: &tokenValConf}
	validationMiddleware := UserAgentValidationMiddleware{Config: &userAgentValConf}
	privilegeMiddleware := AdminPrivilegeMiddleware{}
	
	return Middleware{
		TokenValidation: middleware.Handle(),
		UserAgent:       validationMiddleware.Handle(),
		AdminPrivilege:  privilegeMiddleware.Handle(),
	}
}

type Middleware struct {
	TokenValidation gin.HandlerFunc
	UserAgent       gin.HandlerFunc
	AdminPrivilege  gin.HandlerFunc
}
