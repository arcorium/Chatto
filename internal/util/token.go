package util

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(claims jwt.Claims, duration time.Duration, secretKey string) (string, error) {
	// Set standard claims
	times := time.Now()
	mapClaims := claims.(jwt.MapClaims)
	mapClaims["iat"] = times.Unix()
	mapClaims["nbf"] = times.Unix()
	mapClaims["exp"] = times.Add(duration).Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(secretKey))

	return token, err
}

func ParseToken(tokenString string, validateClaims bool, keyfunc jwt.Keyfunc) (*jwt.Token, error) {
	parser := jwt.Parser{
		SkipClaimsValidation: !validateClaims,
	}
	return parser.Parse(tokenString, keyfunc)
}
func ValidateToken(tokenString string, keyfunc jwt.Keyfunc) error {
	_, err := ParseToken(tokenString, true, keyfunc)
	return err
}
