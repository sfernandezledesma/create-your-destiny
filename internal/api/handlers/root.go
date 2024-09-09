package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
	"github.com/sfernandezledesma/create-your-destiny/internal/models"
)

func RootHandler(c *gin.Context) {
	var data models.DataHome
	data.AllGamesData = cache.GetAllGamesData()
	username := auth.GetUsernameFromContext(c)
	if username != "" {
		data.Username = username
		data.UserGamesData = cache.GetUserGamesData(username)
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func BadRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notFoundPage", nil)
}
