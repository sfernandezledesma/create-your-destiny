package main

import (
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sfernandezledesma/create-your-destiny/internal/api"
	"github.com/sfernandezledesma/create-your-destiny/internal/database"
)

func main() {
	database.InitDB()
	r := api.NewRouter()
	r.Run(":8080")
}
