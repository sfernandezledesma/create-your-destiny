package game

import "github.com/sfernandezledesma/create-your-destiny/internal/utils"

type Game struct {
	Name        string
	Description string
	Scenes      map[utils.Nat]Scene
}

type Scene struct {
	Text  string
	Paths []Path
}

type Path struct {
	Text        string
	Destination utils.Nat
}

type DataHome struct {
	Username  string
	UserGames []string
	AllGames  []string
}

type DataCurrentGame struct {
	Name  string
	Scene Scene
}
