package main

import (
	"fmt"
	"os"
)

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
