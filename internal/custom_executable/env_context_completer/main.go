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

	w1, w2, w3 := trimSlot(wantArg1), trimSlot(wantArg2), trimSlot(wantArg3)
	var bad bool
	if os.Args[1] != w1 {
		completerErr("Expected argv[1] to be '%s' got '%s'\n", w1, os.Args[1])
		bad = true
	}
	if os.Args[2] != w2 {
		completerErr("Expected argv[2] to be '%s' got '%s'\n", w2, os.Args[2])
		bad = true
	}
	if os.Args[3] != w3 {
		completerErr("Expected argv[3] to be '%s' got '%s'\n", w3, os.Args[3])
		bad = true
	}

	eln := trimSlot(envLineVar)
	epn := trimSlot(envPointVar)
	wantL := trimSlot(wantLine)
	wantP := trimSlot(wantPoint)

	gotLine := os.Getenv(eln)
	if gotLine != wantL {
		completerErr("Expected %s to be '%s' got '%s'\n", eln, wantL, gotLine)
		bad = true
	}

	gotPoint := os.Getenv(epn)
	if gotPoint != wantP {
		completerErr("Expected %s to be '%s' got '%s'\n", epn, wantP, gotPoint)
		bad = true
	}

	if bad {
		os.Exit(1)
	}

	fmt.Println(trimSlot(wantCompletion))
}
