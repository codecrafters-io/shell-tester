package main

import (
	"fmt"
	"strings"
)

func main() {
	// Patched at copy time (same mechanism as signature_printer); TrimSpace yields the offered word.
	word := "<<RANDOM_1>>"
	fmt.Println(strings.TrimSpace(word))
}
