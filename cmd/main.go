package main

import (
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sfernandezledesma/create-your-destiny/internal/api"
)

func main() {
	api.NewRouter().Run(":8080")
}
