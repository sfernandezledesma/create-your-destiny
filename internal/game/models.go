package game

type DataHome struct {
	Username  string
	UserGames []string
	AllGames  []string
}

type DataCurrentGame struct {
	Name string
	Page Page
}

type Game struct {
	Pages map[string]Page
}

type Page struct {
	Text    string
	Options []Option
}

type Option struct {
	Text        string
	Destination string
}

var AllGames = []string{"ASD", "BASD", "CASD", "OSD", "BOSD", "COSD"}

var GamesByUser = map[string][]string{
	"asd": {"ASD", "BASD", "CASD"},
	"zxc": {"OSD", "BOSD", "COSD"},
}

var Games = map[string]Game{
	"ASD": {Pages: map[string]Page{
		"1": DataPage,
		"2": {
			"This is page 2",
			[]Option{
				{"Go back to page 1", "1"},
			}}}},
}

var DataPage = Page{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	[]Option{
		{"Go to page 2", "2"},
		{"Go back to page 1", "1"},
		{"Go back to page 3", "3"},
	},
}
