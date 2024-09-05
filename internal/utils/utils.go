package utils

import (
	"log"
	"strconv"
)

type Nat uint32

func ExitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func StringToNat(s string) (Nat, error) {
	nat64, err := strconv.ParseUint(s, 10, 32)
	return Nat(nat64), err
}
