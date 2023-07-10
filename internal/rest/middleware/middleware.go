package middleware

import "chatto/internal/config"

func NewMiddleware(config *config.AppConfig) Middleware {
	tokenValConf := TokenValidationConfig{
		SecretKeyFunc: config.JWTKeyFunc,
		TokenType:     "Bearer",
		SigningType:   config.JWTSigningType,
	}
	userAgentValConf := UserAgentValidationConfig{}

	return Middleware{
		TokenValidation: TokenValidationMiddleware{Config: &tokenValConf},
		UserAgent:       UserAgentValidationMiddleware{Config: &userAgentValConf},
	}
}

type Middleware struct {
	TokenValidation TokenValidationMiddleware
	UserAgent       UserAgentValidationMiddleware
}

const (
	KEY_USER_AGENT = "system"
	KEY_JWT_CLAIMS = "claims"
)
