package util

import (
	"chatto/internal/config"
	"chatto/internal/constant"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(claims jwt.Claims, algo string, duration time.Duration, secretKey string, defaultClaims bool) (string, error) {
	if defaultClaims {
		// Set standard claims
		times := time.Now()
		mapClaims := claims.(jwt.MapClaims)
		mapClaims["iat"] = times.Unix()
		mapClaims["nbf"] = times.Unix()
		mapClaims["exp"] = times.Add(duration).Unix()
	}
	return jwt.NewWithClaims(jwt.GetSigningMethod(algo), claims).SignedString([]byte(secretKey))
}

func CreateAccessToken(claims jwt.Claims, config *config.AppConfig) (string, error) {
	return CreateToken(claims, config.JWTSigningType, time.Duration(config.AccessTokenDuration), config.JWTSecretKey, true)
}

func CreateRefreshToken(claims jwt.Claims, config *config.AppConfig) (string, error) {
	return CreateToken(claims, config.JWTSigningType, time.Duration(config.RefreshTokenDuration), config.JWTSecretKey, true)
}

func ParseToken(tokenString string, keyfunc jwt.Keyfunc) (*jwt.Token, error) {
	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}
	return parser.Parse(tokenString, keyfunc)
}

func ValidateToken(tokenString string, keyfunc jwt.Keyfunc) error {
	token, err := ParseToken(tokenString, keyfunc)
	if err != nil {
		return err
	}

	rawClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	iat_, exist := rawClaims["iat"]
	if !exist {
		return errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}
	exp_, exist := rawClaims["exp"]
	if !exist {
		return errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}

	// Cast
	iat, ok := iat_.(float64)
	if !ok {
		return errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}
	exp, ok := exp_.(float64)
	if !ok {
		return errors.New(constant.MSG_TOKEN_FIELD_INVALID_TYPE)
	}

	if iat > exp {
		return errors.New(constant.MSG_BAD_FORMAT_TOKEN)
	}

	now := float64(time.Now().Unix())
	if exp < now {
		return errors.New(constant.MSG_TOKEN_EXPIRED)
	}

	return nil
}
