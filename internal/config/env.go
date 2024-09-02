package config

import (
	"log"
	"os"
)

var jwtSecret []byte

func GetJWTSecret() []byte {
	if jwtSecret == nil {
		if env, ok := os.LookupEnv("JWT_SECRET"); ok {
			jwtSecret = []byte(env)
		} else {
			log.Fatal("JWT_SECRET environment variable not set")
		}
	}
	return jwtSecret
}
