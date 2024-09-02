package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ExitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckPassword(passwd string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err == nil {
		return true
	} else {
		log.Println(err)
		return false
	}
}
