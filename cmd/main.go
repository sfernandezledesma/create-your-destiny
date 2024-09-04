package main

import "github.com/sfernandezledesma/create-your-destiny/internal/api"

func main() {
	api.NewRouter().Run(":8080")
}
