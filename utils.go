package main

import (
	"fmt"
	"os"
)

const NonUniqueIdError = "non-unique id"

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
