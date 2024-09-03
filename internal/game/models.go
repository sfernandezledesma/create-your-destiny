package game

type DataHome struct {
	Username  string
	UserGames []string
	AllGames  []string
}

type DataCurrentGame struct {
	Name  string
	Scene Scene
}

type Game struct {
	Scenes map[string]Scene
}

type Scene struct {
	Text  string
	Paths []Path
}

type Path struct {
	Text        string
	Destination string
}

var AllGames = []string{"ASD", "BASD", "CASD", "OSD", "BOSD", "COSD"}

var GamesByUser = map[string][]string{
	"asd": {"ASD", "BASD", "CASD"},
	"zxc": {"OSD", "BOSD", "COSD"},
}

var Games = map[string]Game{
	"ASD": {Scenes: map[string]Scene{
		"1": DataScene1,
		"2": {
			"This is scene 2",
			[]Path{
				{"Go back to scene 1", "1"},
			}},
		"3": {
			"This is scene 3",
			[]Path{}},
	},
	},
}

var DataScene1 = Scene{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	[]Path{
		{"Go to scene 2", "2"},
		{"Go to scene 1", "1"},
		{"Go to scene 3", "3"},
	},
}
