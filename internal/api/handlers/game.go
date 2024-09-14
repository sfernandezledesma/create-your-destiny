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

type EditData struct {
	GameData      models.GameData
	GameSceneData models.GameSceneData
}

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

func SaveSettings(c *gin.Context) {
	DB := db.GetDB()
	gameId, err := utils.StringToNat(c.Param("gameId"))
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "errorMessage", err.Error())
		return
	}
	gameName := c.PostForm("gameName")
	description := c.PostForm("description")
	public := false
	if c.PostForm("public") == "on" {
		public = true
	}
	if gameName == "" {
		c.HTML(http.StatusBadRequest, "errorPage", "Name should be filled.")
		return
	}
	_, err = DB.Exec("UPDATE GAME SET NAME=?, DESCRIPTION=?, PUBLIC=? WHERE ID=?", gameName, description, public, gameId)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "errorMessage", err.Error())
	} else {
		gameData := cache.GetGameDataFromId(gameId)
		log.Println(gameData)
		gameData.Name = gameName
		gameData.Description = description
		gameData.Public = public
	}
}

func CreateFormHandler(c *gin.Context) { // username was set in LoggedInMiddleware
	c.HTML(http.StatusOK, "create.html", nil)
}

func EditGameHandler(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId"))
	gameData := cache.GetGameDataFromId(gameId)
	gameSceneData := cache.GetSceneDataFromId(gameId)
	c.HTML(http.StatusOK, "edit.html", EditData{GameData: *gameData, GameSceneData: gameSceneData})
}

func SaveScene(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId"))
	sceneNumber, _ := utils.StringToNat(c.Param("sceneNumber"))
	newDBText := c.PostForm(c.Param("sceneNumber"))
	_, err := db.GetDB().Exec("UPDATE SCENE SET TEXT=? WHERE GAMEID=? AND SCENENUMBER=?", newDBText, gameId, sceneNumber)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errorPage", "Could not update scene.")
	} else {
		cache.UpdateScene(gameId, sceneNumber, newDBText)
	}
}

func NewScene(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId"))
	DB := db.GetDB()
	var maxSceneNumber int
	DB.QueryRow("SELECT COUNT(*) FROM SCENE WHERE SCENE.GAMEID = ?", gameId).Scan(&maxSceneNumber)
	newSceneNumber := utils.Nat(maxSceneNumber + 1)
	_, err := DB.Exec("INSERT INTO SCENE(GAMEID,SCENENUMBER,TEXT) VALUES(?,?,?)", gameId, newSceneNumber, "")
	if err != nil {
		c.HTML(http.StatusBadRequest, "errorPage", "Could not create new scene.")
	} else {
		cache.AddNewScene(gameId, newSceneNumber)
	}
}
