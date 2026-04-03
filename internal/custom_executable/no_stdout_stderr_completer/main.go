package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// No stdout candidates; stderr is patched per line (<<RANDOM_1..3>>), space-padded in the binary.
	l1 := "<<RANDOM_1>>"
	l2 := "<<RANDOM_2>>"
	l3 := "<<RANDOM_3>>"
	fmt.Fprintln(os.Stderr, strings.TrimRight(l1, " "))
	fmt.Fprintln(os.Stderr, strings.TrimRight(l2, " "))
	fmt.Fprintln(os.Stderr, strings.TrimRight(l3, " "))
	time.Sleep(120 * time.Second)
}
