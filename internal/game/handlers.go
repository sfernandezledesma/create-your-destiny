package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/sfernandezledesma/create-your-destiny/internal/config"
)

func BadRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}

func RootHandler(c *gin.Context) {
	var data DataHome
	data.AllGames = AllGames
	var username string
	usernameFromContext, exists := c.Get("username")
	if exists {
		username = usernameFromContext.(string)
	} else {
		tokenString, err := c.Cookie("token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return config.GetJWTSecret(), nil
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
		data.UserGames = GamesByUser[username]
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func PlayHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	pageNumber := c.Param("pageNumber")
	page, ok := Games[gameName].Pages[pageNumber]
	if ok {
		data := DataCurrentGame{Name: gameName, Page: page}
		c.HTML(http.StatusOK, "game.html", data)
	} else {
		BadRouteHandler(c)
	}
}

func EditGameHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	// TODO: send game data to edit
	c.HTML(http.StatusOK, "edit.html", gameName)
}
