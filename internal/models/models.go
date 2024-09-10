package models

import "github.com/sfernandezledesma/create-your-destiny/internal/utils"

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
	Text  string
	Paths []Path
}

func NewScene(text string) Scene {
	return Scene{Text: text, Paths: make([]Path, 2)}
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
	UserGamesData []GameData
	AllGamesData  []GameData
}

type DataCurrentGame struct {
	Id    utils.Nat
	Name  string
	Scene Scene
}
