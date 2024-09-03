package auth

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func CreateTokenString(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
	})
	tokenString, err := token.SignedString(GetJWTSecret())
	return tokenString, err
}

func GetUsernameFromContext(username *string, c *gin.Context) {
	if usernameFromContext, exists := c.Get("username"); exists {
		*username = usernameFromContext.(string)
	} else {
		tokenString, err := c.Cookie("token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return GetJWTSecret(), nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					*username = claims["sub"].(string)
				}
			}
		}
	}
}
