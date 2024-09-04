package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/database"
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
			hash, err := auth.CreateHashFromPassword(password)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Password too long.")
				return
			}
			_, err = database.GetDB().Exec("INSERT INTO USER(NAME, HASH) VALUES(?, ?);", username, hash)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
				return
			}
			RootHandler(c)
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
		var hash []byte
		err := database.GetDB().QueryRow("SELECT HASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
		if err != nil { // this is probably because the user doesn't exist (no rows error)
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", "Username doesn't exist.")
			return
		}
		if auth.CheckPasswordWithHash(password, hash) {
			tokenString, err := auth.CreateTokenString(username)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "login.html", "Server error. Try again later.")
				return
			}
			c.SetCookie("token", tokenString, 34560000, "/", "localhost", false, true)
			c.Set("username", username)
			RootHandler(c)
		} else {
			c.HTML(http.StatusBadRequest, "login.html", "Password is incorrect. Try again.")
		}
	} else {
		c.HTML(http.StatusBadRequest, "login.html", "Fields should not be empty.")
	}
}
