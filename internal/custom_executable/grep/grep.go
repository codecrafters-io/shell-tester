package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: grep PATTERN")
		os.Exit(2) // Standard exit code for grep usage error
	}

	pattern := os.Args[1]
	scanner := bufio.NewScanner(os.Stdin)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			fmt.Println(line)
			found = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "grep: error reading input: %v\n", err)
		os.Exit(2)
	}

	if !found {
		os.Exit(1) // Standard exit code for grep when no lines match
	}

	os.Exit(0)
}
