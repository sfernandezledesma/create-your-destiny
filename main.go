package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "index", dataListGames)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("playHandler")
	book := books[r.PathValue("bookName")]
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "game", book.Pages["1"])
}

func badRouteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("badRouteHandler")
	w.Header().Set("HX-Retarget", "body")
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "notfound", nil)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", rootHandler)
	r.Get("/play/{bookName}", playHandler)
	r.NotFound(badRouteHandler)

	log.Println("Server is starting...")
	check(http.ListenAndServe(":8080", r))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
