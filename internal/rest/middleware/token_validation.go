package middleware

import (
	"errors"
	"net/http"
	"strings"

	"server_client_chat/internal/model"
	"server_client_chat/internal/util"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
)

type TokenValidationConfig struct {
	SecretKeyFunc func(*jwt.Token) (any, error)
	TokenType     string
	SigningType   string
}
type TokenValidationMiddleware struct {
	Config *TokenValidationConfig
}

func (a *TokenValidationMiddleware) Handle() gin.HandlerFunc {
	a.Config.TokenType = "Bearer"

	return func(c *gin.Context) {
		// Get the value
		data := c.GetHeader("Authorization")

		// Parse
		token, err := a.parseToken(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, err.Error(), nil))
			return
		}

		// Validate
		err = a.validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, err.Error(), nil))
			return
		}

		// Set claims on context
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, util.ERR_TOKEN_FORMAT, nil))
			return
		}

		userId := mapClaims["user_id"]
		refreshId := mapClaims["refresh_id"]
		if userId == nil || refreshId == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, util.ERR_TOKEN_FORMAT, nil))
			return
		}

		claims := model.AccessTokenClaims{
			UserId:    userId.(string),
			RefreshId: refreshId.(string),
		}
		c.Set(KEY_JWT_CLAIMS, claims)

		c.Next()
	}
}

func (a *TokenValidationMiddleware) parseToken(fullToken string) (*jwt.Token, error) {
	_, tokenString, err := a.splitHeaderValue(fullToken)
	if err != nil {
		return nil, err
	}

	return util.ParseToken(tokenString, true, a.Config.SecretKeyFunc)
}

func (a *TokenValidationMiddleware) splitHeaderValue(value string) (string, string, error) {
	if len(value) == 0 {
		return "", "", errors.New(util.ERR_NO_ACCESS_TOKEN)
	}

	splitData := strings.Split(value, " ")
	if len(splitData) != 2 {
		return "", "", errors.New(util.ERR_TOKEN_FORMAT)
	}

	return splitData[0], splitData[1], nil
}

func (a *TokenValidationMiddleware) validateToken(token *jwt.Token) error {
	if token.Method.Alg() != a.Config.SigningType {
		return errors.New(util.ERR_TOKEN_FORMAT)
	}
	return nil
}
