package auth

import (
	"time"

	"github.com/OzyKleyton/studio-api/internal/model/user"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var secretKey = viper.GetString("JWT_SECRET")

func GenerateToken(userReq user.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userReq.ID,
		"username": userReq.Username,
		"email":    userReq.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) bool {

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	return err == nil
}
