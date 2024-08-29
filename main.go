package main

import (
	"html/template"
	"log"
	"net/http"
)

type DataListGames struct {
	UserGames   []string
	OthersGames []string
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
	"ASD": {Pages: map[string]Page{"1": dataPage, "2": {}}},
}

var dataPage = Page{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	[]Option{
		{"asdadfgag", "2"},
		{"aslgjalskj", "3"},
	},
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "index", dataListGames)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("playHandler")
	book := books[r.PathValue("bookName")]
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "game", book.Pages["1"])
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /play/{bookName}", playHandler)

	log.Println("Server is starting...")
	check(http.ListenAndServe(":8080", mux))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
