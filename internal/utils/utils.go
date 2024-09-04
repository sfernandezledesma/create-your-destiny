package utils

import (
	"log"
)

func ExitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
