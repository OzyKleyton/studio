package security

import (
	"golang.org/x/crypto/bcrypt"
)

func EncodePassword(password string) ([]byte, error) {
	cost := bcrypt.DefaultCost

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))

	return err == nil
}
