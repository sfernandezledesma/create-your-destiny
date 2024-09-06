package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

func PlayHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	sceneNumber, err := utils.StringToNat(c.Param("sceneNumber"))
	if err != nil {
		BadRouteHandler(c)
		return
	}
	scene, ok := cache.GamesCache[gameName].Scenes[sceneNumber]
	if ok {
		data := game.DataCurrentGame{Name: gameName, Scene: scene}
		c.HTML(http.StatusOK, "game.html", data)
	} else {
		BadRouteHandler(c)
	}
}

func CreateFormHandler(c *gin.Context) { // username was set in LoggedInMiddleware
	// TODO: Everything
	c.HTML(http.StatusOK, "create.html", nil)
}

func EditGameHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	// TODO: send game data to edit
	c.HTML(http.StatusOK, "edit.html", gameName)
}
