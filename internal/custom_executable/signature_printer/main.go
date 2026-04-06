package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	secretCode := "<<RANDOM_1>>"

	fmt.Printf("Program was passed %d args (including program name).\n", len(os.Args))
	fmt.Printf("Arg #0 (program name): %s\n", os.Args[0])

	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Arg #%d: %s\n", i, os.Args[i])
	}

	// Patched value is space-padded to fixed slot width in the binary; trim for stdout (matches tester expectations).
	fmt.Printf("Program Signature: %s\n", strings.TrimRight(secretCode, " "))
}
