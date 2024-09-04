package db

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

var db *sql.DB

func getDB() *sql.DB {
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
	rows, err := getDB().Query("SELECT NAME FROM USER WHERE NAME = ?;", username)
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
	err := getDB().QueryRow("SELECT HASH FROM USER WHERE NAME = ?;", username).Scan(&hash)
	return hash, err
}

func CreateNewUser(username string, hash []byte) error { // Assumes UserExists(username) == false
	_, err := getDB().Exec("INSERT INTO USER(NAME, HASH) VALUES(?, ?);", username, hash)
	return err
}
