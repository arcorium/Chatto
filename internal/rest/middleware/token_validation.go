package middleware

import (
	"chatto/internal/constant"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/util/httputil"

	"errors"
	"net/http"
	"strings"

	"chatto/internal/util"

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
	if len(a.Config.TokenType) == 0 {
		a.Config.TokenType = "Bearer"
	}

	return func(c *gin.Context) {
		// Get the value
		data := c.GetHeader("Authorization")

		_, tokenString, err := a.splitHeaderValue(data)
		if err != nil {
			httputil.ErrorResponse(c, http.StatusUnauthorized, common.NewError(common.AUTH_UNAUTHORIZED, err.Error()))
			c.Abort()
			return
		}
		// Parse
		token, err := a.parseToken(tokenString)
		if err != nil {
			httputil.ErrorResponse(c, http.StatusUnauthorized, common.NewError(common.AUTH_UNAUTHORIZED, constant.MSG_BAD_FORMAT_TOKEN))
			c.Abort()
			return
		}

		// Validate
		err = a.validateToken(token)
		if err != nil {
			httputil.ErrorResponse(c, http.StatusUnauthorized, common.NewError(common.AUTH_TOKEN_NOT_VALIDATED_ERROR, err.Error()))
			c.Abort()
			return
		}

		// Set claims on context
		claims, err := a.GetAccessTokenClaims(token.Claims)
		if err != nil {
			httputil.ErrorResponse(c, http.StatusUnauthorized, common.NewError(common.AUTH_UNAUTHORIZED, err.Error()))
			c.Abort()
			return
		}
		c.Set(constant.KEY_JWT_CLAIMS, &claims)

		c.Next()
	}
}

func (a *TokenValidationMiddleware) parseToken(tokenString string) (*jwt.Token, error) {
	return util.ParseToken(tokenString, a.Config.SecretKeyFunc)
}

func (a *TokenValidationMiddleware) splitHeaderValue(value string) (string, string, error) {
	if len(value) == 0 {
		return "", "", errors.New(constant.MSG_NO_ACCESS_TOKEN)
	}

	splitData := strings.Split(value, " ")
	if len(splitData) != 2 {
		return "", "", errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}

	return splitData[0], splitData[1], nil
}

func (a *TokenValidationMiddleware) validateToken(token *jwt.Token) error {
	if token.Method.Alg() != a.Config.SigningType {
		return errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	return util.ValidateToken(token.Raw, a.Config.SecretKeyFunc)
}

func (a *TokenValidationMiddleware) GetAccessTokenClaims(claims jwt.Claims) (model.AccessTokenClaims, error) {
	result := model.AccessTokenClaims{}
	rawClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return result, errors.New("broken claims")
	}
	userId_, exist := rawClaims["user_id"]
	if !exist {
		return result, errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	refreshId_, exist := rawClaims["refresh_id"]
	if !exist {
		return result, errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	role_, exist := rawClaims["role"]
	if !exist {
		return result, errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	name_, exist := rawClaims["name"]
	if !exist {
		return result, errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}

	userId, ok := userId_.(string)
	if !ok {
		return result, errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}
	refreshId, ok := refreshId_.(string)
	if !ok {
		return result, errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}
	role, ok := role_.(string)
	if !ok {
		return result, errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}
	name, ok := name_.(string)
	if !ok {
		return result, errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}

	result.UserId = userId
	result.Name = name
	result.Role = model.Role(role)
	result.RefreshId = refreshId
	return result, nil
}
