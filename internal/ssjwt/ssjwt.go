package ssjwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/seesawhq/seesaw/config"
)

func CreateLoginToken(config *config.ConfigStruct, userID int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(config.SecretKey))
	return token, err
}

func ParseLoginToken(config *config.ConfigStruct, tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims := token.Claims.(*jwt.StandardClaims)
	tokenInt, err := strconv.ParseInt(claims.Subject, 10, 64)
	return tokenInt, err
}
