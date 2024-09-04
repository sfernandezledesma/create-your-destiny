package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
)

func RootHandler(c *gin.Context) {
	var data game.DataHome
	data.AllGames = game.AllGames
	var username string
	auth.GetUsernameFromContext(&username, c)
	if username != "" {
		data.Username = username
		data.UserGames = game.GamesByUser[username]
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func BadRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}
