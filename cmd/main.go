package main

import "github.com/sfernandezledesma/create-your-destiny/internal/router"

func main() {
	router.NewRouter().Run(":8080")
}
