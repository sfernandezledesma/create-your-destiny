package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"golang.org/x/crypto/bcrypt"
)

type DataListGames struct {
	UserGames   []string
	OthersGames []string
}

type DataCurrentGame struct {
	Name string
	Page Page
}

type Book struct {
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

var dataListGames = DataListGames{
	UserGames:   []string{"ASD", "BASD", "CASD"},
	OthersGames: []string{"OSD", "BOSD", "COSD"},
}

var books = map[string]Book{
	"ASD": {Pages: map[string]Page{
		"1": dataPage,
		"2": {
			"This is page 2",
			[]Option{
				{"Go back to page 1", "1"},
			}}}},
}

var dataPage = Page{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	[]Option{
		{"Go to page 2", "2"},
		{"Go back to page 1", "1"},
	},
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler")
	tmpl := template.Must(template.ParseFiles("assets/index.html"))
	tmpl.ExecuteTemplate(w, "index", dataListGames)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("playHandler")
	bookName := r.PathValue("bookName")
	pageNumber := r.PathValue("pageNumber")
	page := books[bookName].Pages[pageNumber]
	data := DataCurrentGame{Name: bookName, Page: page}
	tmpl := template.Must(template.ParseFiles("assets/index.html"))
	tmpl.ExecuteTemplate(w, "game", data)
}

func badRouteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("badRouteHandler")
	w.Header().Set("HX-Retarget", "body")
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("assets/index.html"))
	tmpl.ExecuteTemplate(w, "notfound", nil)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", rootHandler)
	r.Get("/play/{bookName}/{pageNumber}", playHandler)
	r.NotFound(badRouteHandler)

	// Testing SQLite
	db, err := sql.Open("sqlite3", "app.db")
	check(err)
	rows, err := db.Query("SELECT * FROM USER;")
	check(err)
	defer rows.Close()
	for rows.Next() {
		var name, h string
		check(rows.Scan(&name, &h))
		log.Println(name, h)
		checkPassword("asd123", h)
	}

	// Testing bcrypt
	passwd := "HelloWorld!"
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	check(err)
	checkPassword(passwd, string(hash))

	log.Println("Server is starting...")
	check(http.ListenAndServe(":8080", r))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkPassword(passwd string, hash string) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err == nil {
		log.Println("Password and hash comparison successful!")
	} else {
		log.Fatal(err)
	}
}
