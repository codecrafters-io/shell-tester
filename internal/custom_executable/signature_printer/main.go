package main

import (
	"fmt"
	"os"
)

// This variable will be set at build time.
var secretCode string

func main() {
	var randomCode string
	if len(os.Args) > 1 {
		randomCode = os.Args[1]
	} else {
		randomCode = "<no argument provided>"
	}

	fmt.Printf("Program was passed %d args (including program name).\n", len(os.Args))
	fmt.Printf("Arg #0 (program name): %s\n", os.Args[0])
	fmt.Printf("Arg #1: %s\n", randomCode)
	fmt.Printf("Program Signature: %s\n", secretCode)
}
