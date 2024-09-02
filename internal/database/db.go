package database

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
