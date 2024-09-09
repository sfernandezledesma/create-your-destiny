package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
	"github.com/sfernandezledesma/create-your-destiny/internal/db"
	"github.com/sfernandezledesma/create-your-destiny/internal/models"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

func PlayHandler(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId")) // FIXME: Handle error
	gameData := cache.GetGameDataFromId(gameId)
	gameName := gameData.Name
	sceneNumber, err := utils.StringToNat(c.Param("sceneNumber"))
	if err != nil {
		BadRouteHandler(c)
		return
	}
	scene, ok := cache.GetSceneDataFromId(gameId).GetScene(sceneNumber)
	if ok {
		data := models.DataCurrentGame{Id: gameId, Name: gameName, Scene: scene}
		c.HTML(http.StatusOK, "game.html", data)
	} else {
		BadRouteHandler(c)
	}
}

func CreateGameHandler(c *gin.Context) { // username was set in LoggedInMiddleware
	username := auth.GetUsernameFromContext(c) // This shouldn't be "" because user is logged in
	gameName := c.PostForm("gameName")
	description := c.PostForm("description")
	public := false
	if c.PostForm("public") == "on" {
		public = true
	}
	if gameName == "" {
		c.HTML(http.StatusBadRequest, "create.html", "Fill the form.")
		return
	}
	if err := db.CreateNewGame(gameName, username, description, public); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "create.html", err.Error())
		return
	}
	RootHandler(c)
}

func CreateFormHandler(c *gin.Context) { // username was set in LoggedInMiddleware
	c.HTML(http.StatusOK, "create.html", nil)
}

func EditGameHandler(c *gin.Context) {
	gameId := c.Param("gameId")
	// TODO: send game data to edit
	c.HTML(http.StatusOK, "edit.html", gameId)
}
