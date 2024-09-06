package cache

import (
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

type Cache struct {
}

var AllGames = []string{"ASD", "BASD", "CASD", "OSD", "BOSD", "COSD"}

var GamesByUser = map[string][]string{
	"asd": {"ASD", "BASD", "CASD"},
	"zxc": {"OSD", "BOSD", "COSD"},
}

var GamesCache = map[string]game.Game{
	"ASD": {Scenes: map[utils.Nat]game.Scene{
		1: dataScene1,
		2: {
			Text: "This is scene 2",
			Paths: []game.Path{
				{Text: "Go back to scene 1", Destination: 1},
			}},
		3: {
			Text:  "This is scene 3",
			Paths: []game.Path{}},
	},
	},
}

var dataScene1 = game.Scene{
	Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	Paths: []game.Path{
		{Text: "Go to scene 2", Destination: 2},
		{Text: "Go to scene 1", Destination: 1},
		{Text: "Go to scene 3", Destination: 3},
	},
}
