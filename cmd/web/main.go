package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		{"Go back to page 3", "3"},
	},
}

func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", dataListGames)
}

func playHandler(c *gin.Context) {
	bookName := c.Param("bookName")
	pageNumber := c.Param("pageNumber")
	page, ok := books[bookName].Pages[pageNumber]
	if ok {
		data := DataCurrentGame{Name: bookName, Page: page}
		c.HTML(http.StatusOK, "game.html", data)
	} else {
		badRouteHandler(c)
	}
}

func badRouteHandler(c *gin.Context) {
	c.Header("HX-Retarget", "body")
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}

func registerFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func registerHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" { // FIXME: need better validation
		log.Println(username, password)
		// TODO: add [username, hash] to DB
		rootHandler(c)
	} else {
		c.HTML(http.StatusBadRequest, "register.html", "An error occurred. Try again.") // TODO: send better errors)
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", rootHandler)
	r.GET("/play/:bookName/:pageNumber", playHandler)
	r.GET("/register", registerFormHandler)
	r.POST("/register", registerHandler)
	r.NoMethod(badRouteHandler)
	r.NoRoute(badRouteHandler)

	// Testing SQLite
	db, err := sql.Open("sqlite3", "app.db")
	checkError(err)
	rows, err := db.Query("SELECT * FROM USER;")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		var name, h string
		checkError(rows.Scan(&name, &h))
		log.Println(name, h)
		checkPassword("asd123", h)
	}

	// Testing bcrypt
	passwd := "HelloWorld!"
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	checkError(err)
	checkPassword(passwd, string(hash))

	r.Run(":8080")
}

func checkError(err error) {
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
