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

// Modulus function - Go uses the % operator for remainder (which can return negative results)
func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}
