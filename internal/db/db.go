package db

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

var db *sql.DB

func GetDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", "app.db")
		utils.ExitIfError(err)
	}
	return db
}

func UserExists(username string) (bool, error) {
	var exists bool = false
	var err error = nil
	rows, err := GetDB().Query("SELECT NAME FROM USER WHERE NAME = ?;", username)
	if err == nil {
		defer rows.Close()
		if rows.Next() { // username already exists
			exists = true
		}
	}
	return exists, err
}

func GetUserHash(username string) ([]byte, error) {
	var hash []byte
	err := GetDB().QueryRow("SELECT PASSWORDHASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
	return hash, err
}

func CreateNewUser(username string, hash []byte) error { // Assumes UserExists(username) == false
	_, err := GetDB().Exec("INSERT INTO USER(NAME, PASSWORDHASH) VALUES(?, ?);", username, hash)
	return err
}

func CreateNewGame(gameName string, author string, description string, public bool) error {
	rows, err := GetDB().Query("SELECT NAME FROM GAME WHERE NAME = ?", gameName)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() { // game with that name already exists
		return utils.NewError("Game already exists.")
	}
	_, err = GetDB().Exec("INSERT INTO GAME(NAME, AUTHOR, DESCRIPTION, PUBLIC) VALUES(?, ?, ?, ?)", gameName, author, description, public)
	return err
}
