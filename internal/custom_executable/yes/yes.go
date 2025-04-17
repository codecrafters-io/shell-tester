package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Default output is 'y' if no arguments provided
	output := "y"

	// If arguments are provided, use them as the output string
	if len(os.Args) > 1 {
		output = strings.Join(os.Args[1:], " ")
	}

	// Continuously print the output until the program is terminated
	for {
		fmt.Println(output)
	}
}
