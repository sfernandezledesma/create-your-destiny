package database

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

var SecretKey []byte = []byte("gransecreto")
var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "app.db")
	utils.ExitIfError(err)
}
