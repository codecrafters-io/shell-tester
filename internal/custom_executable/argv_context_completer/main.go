package main

import (
	"fmt"
	"os"
)

// Hardcoded expected invocation (PA6): git remote set<TAB> → set-url
const (
	wantArg1       = "git"
	wantArg2       = "set"
	wantArg3       = "remote"
	completionLine = "set-url"
)

func main() {
	n := len(os.Args) - 1 // number of args after program name
	if n < 3 {
		fmt.Fprintf(os.Stderr, "\nExpected argv[1] thru argv[3], only found up to argv[%d]\n", len(os.Args)-1)
		os.Exit(1)
	}
	if n > 3 {
		fmt.Fprintf(os.Stderr, "\nExpected argv[1] thru argv[3] only, got %d argument(s) after program name\n", n)
		os.Exit(1)
	}

	if os.Args[1] != wantArg1 {
		fmt.Fprintf(os.Stderr, "\nargv[1] mismatch: expected %q, got %q\n", wantArg1, os.Args[1])
		os.Exit(1)
	}
	if os.Args[2] != wantArg2 {
		fmt.Fprintf(os.Stderr, "\nargv[2] mismatch: expected %q, got %q\n", wantArg2, os.Args[2])
		os.Exit(1)
	}
	if os.Args[3] != wantArg3 {
		fmt.Fprintf(os.Stderr, "\nargv[3] mismatch: expected %q, got %q\n", wantArg3, os.Args[3])
		os.Exit(1)
	}

	fmt.Println(completionLine)
}
