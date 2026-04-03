package main

import (
	"fmt"
	"os"
	"strings"
)

// Patched at copy time: slots 1–2 name the env vars to read; 3–5 = expected argv; 6–7 = expected COMP_LINE / COMP_POINT.
var (
	envLineVar     = "<<RANDOM_1>>"
	envPointVar    = "<<RANDOM_2>>"
	wantArg1       = "<<RANDOM_3>>"
	wantArg2       = "<<RANDOM_4>>"
	wantArg3       = "<<RANDOM_5>>"
	wantLine       = "<<RANDOM_6>>"
	wantPoint      = "<<RANDOM_7>>"
	wantCompletion = "<<RANDOM_8>>"
)

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

	w1, w2, w3 := trimSlot(wantArg1), trimSlot(wantArg2), trimSlot(wantArg3)
	if os.Args[1] != w1 {
		fmt.Fprintf(os.Stderr, "\nargv[1] mismatch: expected %q, got %q\n", w1, os.Args[1])
		os.Exit(1)
	}
	if os.Args[2] != w2 {
		fmt.Fprintf(os.Stderr, "\nargv[2] mismatch: expected %q, got %q\n", w2, os.Args[2])
		os.Exit(1)
	}
	if os.Args[3] != w3 {
		fmt.Fprintf(os.Stderr, "\nargv[3] mismatch: expected %q, got %q\n", w3, os.Args[3])
		os.Exit(1)
	}

	eln := trimSlot(envLineVar)
	epn := trimSlot(envPointVar)
	wantL := trimSlot(wantLine)
	wantP := trimSlot(wantPoint)

	gotLine := os.Getenv(eln)
	if gotLine == "" && wantL != "" {
		fmt.Fprintf(os.Stderr, "\nenvironment variable %q is unset or empty (expected %q)\n", eln, wantL)
		os.Exit(1)
	}
	if gotLine != wantL {
		fmt.Fprintf(os.Stderr, "\n%s mismatch: expected %q, got %q\n", eln, wantL, gotLine)
		os.Exit(1)
	}

	gotPoint := os.Getenv(epn)
	if gotPoint == "" && wantP != "" {
		fmt.Fprintf(os.Stderr, "\nenvironment variable %q is unset or empty (expected %q)\n", epn, wantP)
		os.Exit(1)
	}
	if gotPoint != wantP {
		fmt.Fprintf(os.Stderr, "\n%s mismatch: expected %q, got %q\n", epn, wantP, gotPoint)
		os.Exit(1)
	}

	fmt.Println(trimSlot(wantCompletion))
}
