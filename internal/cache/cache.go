package cache

import (
	"github.com/sfernandezledesma/create-your-destiny/internal/db"
	"github.com/sfernandezledesma/create-your-destiny/internal/models"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

var gameDataById map[utils.Nat]models.GameData
var sceneDataById map[utils.Nat]*models.GameSceneData
var gamesByUser map[string][]models.GameData

func InitCache() { // Should be called on server start
	gameDataById = make(map[utils.Nat]models.GameData)
	sceneDataById = make(map[utils.Nat]*models.GameSceneData)
	gamesByUser = make(map[string][]models.GameData)
	DB := db.GetDB()
	var id, sceneNumber, originScene, destScene utils.Nat
	var name, author, description, text string
	var public bool
	rows, err := DB.Query("SELECT ID, NAME, AUTHOR, DESCRIPTION, PUBLIC FROM GAME")
	if err == nil {
		for rows.Next() {
			if err := rows.Scan(&id, &name, &author, &description, &public); err == nil {
				gameData := models.NewGameData(id, name, author, description)
				if public {
					gameDataById[id] = gameData
				}
				sceneDataById[id] = new(models.GameSceneData)
				*sceneDataById[id] = models.NewGameSceneData(id)
				if gamesByUser[author] == nil {
					gamesByUser[author] = make([]models.GameData, 0, 4)
				}
				gamesByUser[author] = append(gamesByUser[author], gameData)
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
				newScene := new(models.Scene)
				newScene.Text = text
				sceneDataById[id].Scenes[sceneNumber] = newScene
			}
		}
	}
	rows.Close()
	rows, err = DB.Query("SELECT GAMEID, ORIGINSCENE, DESTSCENE, TEXT FROM PATH")
	if err == nil {
		for rows.Next() {
			if err := rows.Scan(&id, &originScene, &destScene, &text); err == nil {
				sceneDataById[id].Scenes[originScene].Paths = append(sceneDataById[id].Scenes[originScene].Paths, models.NewPath(text, destScene))
			}
		}
	}
}

func GetAllGamesData() []models.GameData {
	gamesData := make([]models.GameData, 0, len(gameDataById))
	for _, v := range gameDataById {
		gamesData = append(gamesData, v)
	}
	return gamesData
}

func GetGameDataFromId(id utils.Nat) models.GameData {
	return gameDataById[id]
}

func GetUserGamesData(username string) []models.GameData {
	return gamesByUser[username]
}

func GetSceneDataFromId(gameId utils.Nat) models.GameSceneData {
	return *sceneDataById[gameId]
}

// var GamesCache = map[string]models.Game{
// 	"ASD": {Scenes: map[utils.Nat]models.Scene{
// 		1: dataScene1,
// 		2: {
// 			Text: "This is scene 2",
// 			Paths: []models.Path{
// 				{Text: "Go back to scene 1", Destination: 1},
// 			}},
// 		3: {
// 			Text:  "This is scene 3",
// 			Paths: []models.Path{}},
// 	},
// 	},
// }

// var dataScene1 = models.Scene{
// 	Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
// 	Paths: []models.Path{
// 		{Text: "Go to scene 2", Destination: 2},
// 		{Text: "Go to scene 1", Destination: 1},
// 		{Text: "Go to scene 3", Destination: 3},
// 	},
// }
