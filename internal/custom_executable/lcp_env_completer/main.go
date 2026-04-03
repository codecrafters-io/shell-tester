package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Patched at copy time: env var names for COMP_LINE / COMP_POINT, expected argv[1] (command) and argv[3] (previous word).
var (
	envLineVar  = "<<RANDOM_1>>"
	envPointVar = "<<RANDOM_2>>"
	wantArg1    = "<<RANDOM_3>>"
	wantArg3    = "<<RANDOM_4>>"
)

var allCandidates = []string{"checkout", "cherry-pick"}

var completerErrHeader bool

func completerErr(format string, a ...any) {
	if !completerErrHeader {
		fmt.Fprintf(os.Stderr, "\nError from the completer script:\n")
		completerErrHeader = true
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

func trimSlot(s string) string {
	return strings.TrimRight(s, " ")
}

func main() {
	n := len(os.Args) - 1
	if n < 3 {
		completerErr("Expected argv[1] thru argv[3], only found up to argv[%d]\n", len(os.Args)-1)
		os.Exit(1)
	}
	if n > 3 {
		completerErr("Expected argv[1] thru argv[3] only, got %d argument(s) after program name\n", n)
		os.Exit(1)
	}

	w1, w3 := trimSlot(wantArg1), trimSlot(wantArg3)
	var bad bool
	if os.Args[1] != w1 {
		completerErr("Expected argv[1] to be '%s' got '%s'\n", w1, os.Args[1])
		bad = true
	}
	if os.Args[3] != w3 {
		completerErr("Expected argv[3] to be '%s' got '%s'\n", w3, os.Args[3])
		bad = true
	}

	eln := trimSlot(envLineVar)
	epn := trimSlot(envPointVar)
	gotLine := os.Getenv(eln)
	gotPointStr := os.Getenv(epn)

	if gotLine == "" {
		completerErr("Expected %s to be non-empty got ''\n", eln)
		bad = true
	}

	wantPoint := fmt.Sprintf("%d", len(gotLine))
	if gotPointStr != wantPoint {
		completerErr("Expected %s to be '%s' got '%s'\n", epn, wantPoint, gotPointStr)
		bad = true
	}

	if gotLine != "" && !strings.HasPrefix(gotLine, w1+" ") {
		completerErr("Expected %s to start with '%s' got '%s'\n", eln, w1+" ", gotLine)
		bad = true
	}

	if bad {
		os.Exit(1)
	}

	prefix := os.Args[2]
	var matches []string
	for _, w := range allCandidates {
		if strings.HasPrefix(w, prefix) {
			matches = append(matches, w)
		}
	}
	sort.Strings(matches)
	for _, w := range matches {
		fmt.Println(w)
	}
}
