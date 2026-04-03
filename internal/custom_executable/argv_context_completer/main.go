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

var completerErrHeader bool

func completerErr(format string, a ...any) {
	if !completerErrHeader {
		fmt.Fprintf(os.Stderr, "\nError from the completer script:\n")
		completerErrHeader = true
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

func main() {
	n := len(os.Args) - 1 // number of args after program name
	if n < 3 {
		completerErr("Expected argv[1] thru argv[3], only found up to argv[%d]\n", len(os.Args)-1)
		os.Exit(1)
	}
	if n > 3 {
		completerErr("Expected argv[1] thru argv[3] only, got %d argument(s) after program name\n", n)
		os.Exit(1)
	}

	var bad bool
	if os.Args[1] != wantArg1 {
		completerErr("Expected argv[1] to be '%s' got '%s'\n", wantArg1, os.Args[1])
		bad = true
	}
	if os.Args[2] != wantArg2 {
		completerErr("Expected argv[2] to be '%s' got '%s'\n", wantArg2, os.Args[2])
		bad = true
	}
	if os.Args[3] != wantArg3 {
		completerErr("Expected argv[3] to be '%s' got '%s'\n", wantArg3, os.Args[3])
		bad = true
	}
	if bad {
		os.Exit(1)
	}

	fmt.Println(completionLine)
}
