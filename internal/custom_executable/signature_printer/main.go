package main

import (
	"fmt"
	"os"
)

func main() {
	secretCode := "<<RANDOM>>"

	fmt.Printf("Program was passed %d args (including program name).\n", len(os.Args))
	fmt.Printf("Arg #0 (program name): %s\n", os.Args[0])

	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Arg #%d: %s\n", i, os.Args[i])
	}

	fmt.Printf("Program Signature: %s\n", secretCode)
}
