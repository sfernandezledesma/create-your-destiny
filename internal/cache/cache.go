package cache

import (
	"github.com/sfernandezledesma/create-your-destiny/internal/db"
	"github.com/sfernandezledesma/create-your-destiny/internal/models"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

var gameDataById map[utils.Nat]*models.GameData
var sceneDataById map[utils.Nat]*models.GameSceneData
var gamesByUser map[string][]*models.GameData

func InitCache() { // Should be called on server start
	gameDataById = make(map[utils.Nat]*models.GameData)
	sceneDataById = make(map[utils.Nat]*models.GameSceneData)
	gamesByUser = make(map[string][]*models.GameData)
	DB := db.GetDB()
	var id, sceneNumber utils.Nat
	var name, author, description, text string
	var public bool
	rows, err := DB.Query("SELECT ID, NAME, AUTHOR, DESCRIPTION, PUBLIC FROM GAME")
	if err == nil {
		for rows.Next() {
			if err := rows.Scan(&id, &name, &author, &description, &public); err == nil {
				gameData := models.NewGameData(id, name, author, description, public)
				gameDataById[id] = &gameData
				sceneDataById[id] = new(models.GameSceneData)
				*sceneDataById[id] = models.NewGameSceneData(id)
				if gamesByUser[author] == nil {
					gamesByUser[author] = make([]*models.GameData, 0, 4)
				}
				gamesByUser[author] = append(gamesByUser[author], &gameData)
			}
		}
	}
	rows.Close()
	rows, err = DB.Query("SELECT GAMEID, SCENENUMBER, TEXT FROM SCENE")
	if err == nil {
		for rows.Next() {
			if err := rows.Scan(&id, &sceneNumber, &text); err == nil {
				if sceneDataById[id].Scenes == nil {
					sceneDataById[id].Scenes = make(map[utils.Nat]*models.Scene)
				}
				sceneDataById[id].Scenes[sceneNumber] = models.NewSceneFromCode(text)
			}
		}
	}
	rows.Close()
}

func GetAllPublicGamesData() []*models.GameData {
	gamesData := make([]*models.GameData, 0, len(gameDataById))
	for _, gameData := range gameDataById {
		if gameData.Public {
			gamesData = append(gamesData, gameData)
		}
	}
	return gamesData
}

func GetGameDataFromId(id utils.Nat) *models.GameData {
	return gameDataById[id]
}

func GetUserGamesData(username string) []*models.GameData {
	return gamesByUser[username]
}

func GetSceneDataFromId(gameId utils.Nat) models.GameSceneData {
	return *sceneDataById[gameId]
}

func AddNewScene(gameId utils.Nat, newSceneNumber utils.Nat) {
	newScene := new(models.Scene)
	newScene.Text = ""
	sceneDataById[gameId].Scenes[newSceneNumber] = newScene
}

func UpdateScene(gameId utils.Nat, sceneNumber utils.Nat, newDBText string) {
	sceneDataById[gameId].Scenes[sceneNumber] = models.NewSceneFromCode(newDBText)
}
