package main

import (
	"github.com/sfernandezledesma/create-your-destiny/internal/api"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
)

func main() {
	cache.InitCache()
	api.NewRouter().Run(":8080")
}
