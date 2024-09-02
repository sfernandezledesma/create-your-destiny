package auth

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sfernandezledesma/create-your-destiny/internal/config"
	"github.com/sfernandezledesma/create-your-destiny/internal/database"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		rows, err := database.GetDB().Query("SELECT NAME FROM USER WHERE NAME = ?;", username)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
			return
		}
		defer rows.Close()
		if rows.Next() { // username already exists
			c.HTML(http.StatusBadRequest, "register.html", "Username already exists.")
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Password too long.")
				return
			}
			utils.CheckPassword(password, string(hash))
			_, err = database.GetDB().Exec("INSERT INTO USER(NAME, HASH) VALUES(?, ?);", username, hash)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
				return
			}
			game.RootHandler(c)
		}
	} else {
		c.HTML(http.StatusBadRequest, "register.html", "Fields should not be empty.")
	}
}

func LoginFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		var hash string
		err := database.GetDB().QueryRow("SELECT HASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
		if err != nil { // this is probably because the user doesn't exist (no rows error)
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", "Username doesn't exist.")
			return
		}
		if utils.CheckPassword(password, hash) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": username,
			})
			tokenString, err := token.SignedString(config.GetJWTSecret())
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "login.html", "Server error. Try again later.")
				return
			}
			c.SetCookie("token", tokenString, 34560000, "/", "localhost", false, true)
			c.Set("username", username)
			game.RootHandler(c)
		} else {
			c.HTML(http.StatusBadRequest, "login.html", "Password is incorrect. Try again.")
		}
	} else {
		c.HTML(http.StatusBadRequest, "login.html", "Fields should not be empty.")
	}
}

func GameOwnerMiddleware(c *gin.Context) {
	gameName := c.Param("gameName")

	// Check if the user is logged in and retrieve username
	var username string
	tokenString, err := c.Cookie("token")
	if err == nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				username = claims["sub"].(string)
				c.Set("username", username)
			}
		}
	}

	if username == "" {
		c.HTML(http.StatusUnauthorized, "errorPage", "Unauthorized")
		c.Abort()
		return
	}

	// Check if the user is the owner of the game, gameName is unique
	if !slices.Contains(game.GamesByUser[username], gameName) {
		c.HTML(http.StatusForbidden, "errorPage", "Forbidden")
		c.Abort()
		return
	}

	// If everything is fine, proceed to the next handler
	c.Next()
}
