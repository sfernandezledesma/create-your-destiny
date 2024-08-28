package main

import (
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	UserGames   []string
	OthersGames []string
}

var data = Data{
	UserGames:   []string{"ASD", "BASD", "CASD"},
	OthersGames: []string{"OSD", "BOSD", "COSD"},
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", rootHandler)

	check(http.ListenAndServe(":8080", nil))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
