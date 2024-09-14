package models

import (
	"regexp"
	"strings"

	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

type GameData struct {
	Id          utils.Nat
	Name        string
	Author      string
	Description string
	Public      bool
}

func NewGameData(id utils.Nat, name string, author string, description string, public bool) GameData {
	return GameData{Id: id, Name: name, Author: author, Description: description, Public: public}
}

type GameSceneData struct {
	GameId utils.Nat
	Scenes map[utils.Nat]*Scene
}

func (gsd GameSceneData) GetScene(sceneNumber utils.Nat) (Scene, bool) {
	scene, ok := gsd.Scenes[sceneNumber]
	return *scene, ok
}

func NewGameSceneData(gameId utils.Nat) GameSceneData {
	return GameSceneData{GameId: gameId, Scenes: make(map[utils.Nat]*Scene, 16)}
}

type Scene struct {
	DBText string
	Text   string
	Paths  []Path
}

func NewSceneFromCode(dbText string) *Scene { // This also instantiates the Paths and add them to the new scene
	newScene := new(Scene)
	newScene.DBText = dbText
	lines := strings.Split(dbText, "\n")
	re := regexp.MustCompile(`\[(\d+)\] (.+)`)
	sceneText := ""
	for _, line := range lines {
		if match := re.FindSubmatch([]byte(line)); match != nil {
			pathDestination, _ := utils.StringToNat(string(match[1]))
			pathText := string(match[2])
			newScene.addPath(NewPath(pathText, pathDestination))
		} else {
			sceneText += (line + "\n")
		}
	}
	newScene.Text = sceneText
	return newScene
}

func (s *Scene) addPath(path Path) {
	s.Paths = append(s.Paths, path)
}

type Path struct {
	Text        string
	Destination utils.Nat
}

func NewPath(text string, destination utils.Nat) Path {
	return Path{Text: text, Destination: destination}
}

type DataHome struct {
	Username      string
	UserGamesData []*GameData
	AllGamesData  []*GameData
}

type DataCurrentGame struct {
	Id    utils.Nat
	Name  string
	Scene Scene
}
