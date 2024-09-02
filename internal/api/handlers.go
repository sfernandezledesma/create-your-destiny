package api

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/sfernandezledesma/create-your-destiny/internal/database"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

func playHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	pageNumber := c.Param("pageNumber")
	page, ok := game.Games[gameName].Pages[pageNumber]
	if ok {
		data := game.DataCurrentGame{Name: gameName, Page: page}
		c.HTML(http.StatusOK, "game.html", data)
	} else {
		badRouteHandler(c)
	}
}

func badRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}

func registerFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func registerHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		log.Println(username, password)
		rows, err := database.DB.Query("SELECT NAME FROM USER WHERE NAME = ?;", username)
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
			result, err := database.DB.Exec("INSERT INTO USER(NAME, HASH) VALUES(?, ?);", username, hash)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
				return
			}
			log.Println(result)
			rootHandler(c)
		}
	} else {
		c.HTML(http.StatusBadRequest, "register.html", "Fields should not be empty.")
	}
}

func loginFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		log.Println(username, password)
		var hash string
		err := database.DB.QueryRow("SELECT HASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
		if err != nil { // this is probably because the user doesn't exist (no rows error)
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", "Username doesn't exist.")
			return
		}
		if utils.CheckPassword(password, hash) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": username,
			})
			tokenString, err := token.SignedString(database.SecretKey)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "login.html", "Server error. Try again later.")
				return
			}
			c.SetCookie("token", tokenString, 34560000, "/", "localhost", false, true)
			c.Set("username", username)
			rootHandler(c)
		} else {
			c.HTML(http.StatusBadRequest, "login.html", "Password is incorrect. Try again.")
		}
	} else {
		c.HTML(http.StatusBadRequest, "login.html", "Fields should not be empty.")
	}
}

func rootHandler(c *gin.Context) {
	var data game.DataHome
	data.AllGames = game.AllGames
	var username string
	usernameFromContext, exists := c.Get("username")
	if exists {
		username = usernameFromContext.(string)
	} else {
		tokenString, err := c.Cookie("token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return database.SecretKey, nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					username = claims["sub"].(string)
				}
			}
		}
	}
	if username != "" {
		data.Username = username
		data.UserGames = game.GamesByUser[username]
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func gameOwnerMiddleware(c *gin.Context) {
	gameName := c.Param("gameName")

	// Check if the user is logged in and retrieve username
	var username string
	tokenString, err := c.Cookie("token")
	if err == nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return database.SecretKey, nil
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

func editGameHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	// TODO: send game data to edit
	c.HTML(http.StatusOK, "edit.html", gameName)
}
