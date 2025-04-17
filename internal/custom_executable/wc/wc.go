package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type counts struct {
	lines, words, chars int
}

func main() {
	flagSet := flag.NewFlagSet("wc", flag.ContinueOnError)
	// Define flags
	lFlag := flagSet.Bool("l", false, "count lines")
	wFlag := flagSet.Bool("w", false, "count words")
	cFlag := flagSet.Bool("c", false, "count characters")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "wc: invalid option: %v\n", err)
		os.Exit(1)
	}

	// If no flags are provided, show all counts
	showAll := !(*lFlag || *wFlag || *cFlag)
	if showAll {
		*lFlag = true
		*wFlag = true
		*cFlag = true
	}

	fileArgs := flagSet.Args()
	if len(fileArgs) == 0 {
		// If no files provided, read from stdin
		counts := countReader(os.Stdin)
		printCounts(counts, "-", *lFlag, *wFlag, *cFlag)
		return
	}

	// Process each file
	total := counts{}
	exitWithError := false
	for _, file := range fileArgs {
		if !checkIfFileExists(file) {
			fmt.Fprintf(os.Stderr, "wc: %s: No such file or directory\n", file)
			exitWithError = true
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", file, err)
			exitWithError = true
			continue
		}

		counts := countReader(f)
		printCounts(counts, file, *lFlag, *wFlag, *cFlag)

		total.lines += counts.lines
		total.words += counts.words
		total.chars += counts.chars

		f.Close()
	}

	// Print total if multiple files
	if len(fileArgs) > 1 {
		printCounts(total, "total", *lFlag, *wFlag, *cFlag)
	}

	if exitWithError {
		os.Exit(1)
	}
	os.Exit(0)
}

func countReader(r io.Reader) counts {
	var c counts
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		c.lines++
		line := scanner.Text()
		c.words += len(strings.Fields(line))
		c.chars += len(line) + 1 // +1 for newline
	}

	return c
}

func printCounts(c counts, name string, showLines, showWords, showChars bool) {
	if showLines {
		fmt.Printf("%8d", c.lines)
	}
	if showWords {
		fmt.Printf("%8d", c.words)
	}
	if showChars {
		fmt.Printf("%8d", c.chars)
	}
	fmt.Printf(" %s\n", name)
}

func checkIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
