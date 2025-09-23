package auth

import (
	"time"

	"github.com/OzyKleyton/studio-api/internal/model/user"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var secretKey = viper.GetString("JWT_SECRET")

type CustomClaims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userReq user.User) (string, error) {

	claims := CustomClaims{
		UserID:   userReq.ID,
		Username: userReq.Username,
		Email:    userReq.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "studio-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func IsTokenValid(tokenString string) bool {
	_, err := VerifyToken(tokenString)
	return err == nil
}
