package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CreateHashFromPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return hash, err
}

func CheckPasswordWithHash(passwd string, hash []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(passwd)); err == nil {
		return true
	} else {
		log.Println(err)
		return false
	}
}
