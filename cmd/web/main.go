package main

import (
	"database/sql"
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"golang.org/x/crypto/bcrypt"
)

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

var allGames = []string{"ASD", "BASD", "CASD", "OSD", "BOSD", "COSD"}

var gamesByUser = map[string][]string{
	"asd": {"ASD", "BASD", "CASD"},
	"zxc": {"OSD", "BOSD", "COSD"},
}

var games = map[string]Game{
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

func playHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	pageNumber := c.Param("pageNumber")
	page, ok := games[gameName].Pages[pageNumber]
	if ok {
		data := DataCurrentGame{Name: gameName, Page: page}
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
	if username != "" && password != "" {
		log.Println(username, password)
		rows, err := db.Query("SELECT NAME FROM USER WHERE NAME = ?;", username)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
			return
		}
		defer rows.Close()
		if rows.Next() { // username already exists
			c.HTML(http.StatusBadRequest, "register.html", "Username already exists.")
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Password too long.")
				return
			}
			checkPassword(password, string(hash))
			result, err := db.Exec("INSERT INTO USER(NAME, HASH) VALUES(?, ?);", username, hash)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "register.html", "Database error. Try again later.")
				return
			}
			log.Println(result)
			rootHandler(c)
		}
	} else {
		c.HTML(http.StatusBadRequest, "register.html", "Fields should not be empty.")
	}
}

func loginFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		log.Println(username, password)
		var hash string
		err := db.QueryRow("SELECT HASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "login.html", "Database error. Try again later.")
			return
		}
		if checkPassword(password, hash) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": username,
			})
			tokenString, err := token.SignedString(secretkey)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "login.html", "Server error. Try again later.")
				return
			}
			c.SetCookie("token", tokenString, 34560000, "/", "localhost", false, true)
			c.Set("username", username)
			rootHandler(c)
		} else {
			c.HTML(http.StatusBadRequest, "login.html", "Password is incorrect. Try again.")
		}
	} else {
		c.HTML(http.StatusBadRequest, "login.html", "Fields should not be empty.")
	}
}

func rootHandler(c *gin.Context) {
	var data DataHome
	data.AllGames = allGames
	var username string
	usernameFromContext, exists := c.Get("username")
	if exists {
		username = usernameFromContext.(string)
	} else {
		tokenString, err := c.Cookie("token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretkey, nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					username = claims["sub"].(string)
				}
			}
		}
	}
	if username != "" {
		data.Username = username
		data.UserGames = gamesByUser[username]
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func gameOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		gameName := c.Param("gameName")

		// Check if the user is logged in and retrieve username
		var username string
		tokenString, err := c.Cookie("token")
		if err == nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretkey, nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					username = claims["sub"].(string)
					c.Set("username", username)
				}
			}
		}

		if username == "" {
			c.HTML(http.StatusUnauthorized, "errorPage", "Unauthorized")
			c.Abort()
			return
		}

		// Check if the user is the owner of the game, gameName is unique
		if !slices.Contains(gamesByUser[username], gameName) {
			c.HTML(http.StatusForbidden, "errorPage", "Forbidden")
			c.Abort()
			return
		}

		// If everything is fine, proceed to the next handler
		c.Next()
	}
}

func editGameHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	// TODO: send game data to edit
	c.HTML(http.StatusOK, "edit.html", gameName)
}

var db *sql.DB
var secretkey []byte = []byte("gransecreto")

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", rootHandler)
	r.GET("/register", registerFormHandler)
	r.POST("/register", registerHandler)
	r.GET("/login", loginFormHandler)
	r.POST("/login", loginHandler)
	r.GET("/play/:gameName/:pageNumber", playHandler)
	r.GET("/edit/:gameName", gameOwnerMiddleware(), editGameHandler)
	r.NoMethod(badRouteHandler)
	r.NoRoute(badRouteHandler)

	var err error
	db, err = sql.Open("sqlite3", "app.db")
	exitIfError(err)
	// rows, err := db.Query("SELECT * FROM USER;")
	// checkError(err)
	// defer rows.Close()
	// for rows.Next() {
	// 	var name, h string
	// 	checkError(rows.Scan(&name, &h))
	// 	log.Println(name, h)
	// }

	r.Run(":8080")
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkPassword(passwd string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err == nil {
		log.Println("Password and hash comparison successful!")
		return true
	} else {
		log.Println(err)
		return false
	}
}
