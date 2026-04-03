package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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

func trimSlot(s string) string {
	return strings.TrimRight(s, " ")
}

func main() {
	n := len(os.Args) - 1
	if n < 3 {
		fmt.Fprintf(os.Stderr, "\nExpected argv[1] thru argv[3], only found up to argv[%d]\n", len(os.Args)-1)
		os.Exit(1)
	}
	if n > 3 {
		fmt.Fprintf(os.Stderr, "\nExpected argv[1] thru argv[3] only, got %d argument(s) after program name\n", n)
		os.Exit(1)
	}

	w1, w3 := trimSlot(wantArg1), trimSlot(wantArg3)
	if os.Args[1] != w1 {
		fmt.Fprintf(os.Stderr, "\nargv[1] mismatch: expected %q, got %q\n", w1, os.Args[1])
		os.Exit(1)
	}
	if os.Args[3] != w3 {
		fmt.Fprintf(os.Stderr, "\nargv[3] mismatch: expected %q, got %q\n", w3, os.Args[3])
		os.Exit(1)
	}

	eln := trimSlot(envLineVar)
	epn := trimSlot(envPointVar)
	gotLine := os.Getenv(eln)
	gotPointStr := os.Getenv(epn)
	if gotLine == "" {
		fmt.Fprintf(os.Stderr, "\nenvironment variable %q is unset or empty\n", eln)
		os.Exit(1)
	}
	if gotPointStr == "" {
		fmt.Fprintf(os.Stderr, "\nenvironment variable %q is unset or empty\n", epn)
		os.Exit(1)
	}
	point, err := strconv.Atoi(gotPointStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s is not a valid integer: %q\n", epn, gotPointStr)
		os.Exit(1)
	}
	if point != len(gotLine) {
		fmt.Fprintf(os.Stderr, "\n%s (%d) must equal byte length of %s (%d) when the cursor is at end-of-line\n", epn, point, eln, len(gotLine))
		os.Exit(1)
	}
	if !strings.HasPrefix(gotLine, w1+" ") {
		fmt.Fprintf(os.Stderr, "\n%s must start with %q followed by a space, got %q\n", eln, w1, gotLine)
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
