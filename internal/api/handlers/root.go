package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
)

func RootHandler(c *gin.Context) {
	var data game.DataHome
	data.AllGames = cache.AllGames
	username := auth.GetUsernameFromContext(c)
	if username != "" {
		data.Username = username
		data.UserGames = cache.GamesByUser[username]
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func BadRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notFoundPage", nil)
}
