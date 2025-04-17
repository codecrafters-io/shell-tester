package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// byteCounter implements io.Writer to count bytes written to it.
type byteCounter struct {
	count int64
}

func (bc *byteCounter) Write(p []byte) (int, error) {
	bc.count += int64(len(p))
	return len(p), nil
}

// counts holds the counts for lines, words, and bytes.
type counts struct {
	lines, words, bytes int64
}

func main() {
	flagSet := flag.NewFlagSet("wc", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard) // Suppress default flag errors

	lFlag := flagSet.Bool("l", false, "count lines")
	wFlag := flagSet.Bool("w", false, "count words")
	cFlag := flagSet.Bool("c", false, "count bytes")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "wc: invalid option: %v\n", err)
		os.Exit(1) // Exit code 1 for bad flags
	}

	showAll := !*lFlag && !*wFlag && !*cFlag
	if showAll {
		*lFlag = true
		*wFlag = true
		*cFlag = true
	}

	fileArgs := flagSet.Args()
	exitCode := 0 // Assume success

	if len(fileArgs) == 0 {
		// Read from stdin
		// Pass flags down to countReader for potential optimization (though not strictly necessary here)
		counts, err := countReader(os.Stdin, *lFlag, *wFlag, *cFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: error reading stdin: %v\n", err)
			os.Exit(1)
		}
		printCounts(counts, "-", *lFlag, *wFlag, *cFlag)
		os.Exit(0)
	}

	total := counts{}

	// Print errors first
	for _, filename := range fileArgs {
		_, err := os.Open(filename)
		if err != nil {
			// Match specific error message for "No such file"
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "wc: %s: open: No such file or directory\n", filename)
			} else {
				fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
			}
			exitCode = 1
		}
	}

	for _, filename := range fileArgs {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}

		// Pass flags to countReader
		fileCounts, readErr := countReader(f, *lFlag, *wFlag, *cFlag)
		closeErr := f.Close() // Close file immediately after reading

		// Prioritize reporting read errors
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "wc: error reading %s: %v\n", filename, readErr)
			exitCode = 1
			// No need to close again, but continue to next file
			continue
		}
		// Report close error if read succeeded but close failed
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "wc: error closing %s: %v\n", filename, closeErr)
			exitCode = 1
			// Still proceed to print counts and add to total, but ensure non-zero exit
		}

		// Print counts for the current file
		displayName := filepath.Base(filename)
		printCounts(fileCounts, displayName, *lFlag, *wFlag, *cFlag)

		// Add to total counts
		total.lines += fileCounts.lines
		total.words += fileCounts.words
		total.bytes += fileCounts.bytes
	}

	// Print total if multiple file arguments were provided
	if len(fileArgs) > 1 {
		printCounts(total, "total", *lFlag, *wFlag, *cFlag)
	}

	os.Exit(exitCode)
}

// countReader now accepts flags to potentially optimize (e.g., skip word count if !wFlag)
func countReader(r io.Reader, countLines, countWords, countBytes bool) (counts, error) {
	var c counts
	var readerForScanner io.Reader = r // Start with the original reader

	// Setup byte counter only if needed
	var counter *byteCounter
	if countBytes {
		counter = &byteCounter{}
		// If we need bytes, tee the original reader
		readerForScanner = io.TeeReader(r, counter)
	}

	// Scanner reads from the (potentially tee'd) reader
	scanner := bufio.NewScanner(readerForScanner)

	for scanner.Scan() {
		// Always count lines if requested or if words are needed (scanner works line by line)
		if countLines || countWords {
			c.lines++
		}
		// Only count words if requested
		if countWords {
			c.words += int64(len(strings.Fields(scanner.Text())))
		}
		// Byte counting happens via the TeeReader automatically if countBytes is true
	}

	if err := scanner.Err(); err != nil {
		return counts{}, err // Return zero counts and the error
	}

	// Assign byte count if it was calculated
	if countBytes && counter != nil {
		c.bytes = counter.count
	}

	if countBytes && !countLines && !countWords {
		// Special case: only count bytes
		// Lines and words remain 0
		return c, nil
	}

	scanner = bufio.NewScanner(readerForScanner) // Re-init scanner with potentially tee'd reader

	for scanner.Scan() {
		c.lines++
		c.words += int64(len(strings.Fields(scanner.Text())))
	}

	if err := scanner.Err(); err != nil {
		return counts{}, err // Return zero counts and the error
	}
	if countBytes && counter != nil { // Assign bytes if counted
		c.bytes = counter.count
	}
	return c, nil
}

// printCounts formats and prints the counts based on active flags.
// Uses %8d for each field, joined by a single space.
func printCounts(c counts, name string, showLines, showWords, showBytes bool) {
	fields := []string{} // Store formatted fields to be printed

	if showLines {
		fields = append(fields, fmt.Sprintf("%8d", c.lines))
	}
	if showWords {
		fields = append(fields, fmt.Sprintf("%8d", c.words))
	}
	if showBytes {
		fields = append(fields, fmt.Sprintf("%8d", c.bytes))
	}

	// Join the count fields with a single space
	output := strings.Join(fields, "")

	// Append the name if it's not the stdin placeholder "-"
	if name != "-" {
		// Add a space before the name only if there were count fields printed
		if len(fields) > 0 {
			output += " " + name
		} else {
			output += name // This line might be problematic if fields is empty.
		}
	}

	// Print the output line only if there's something to print
	// This correctly handles the case where fields is empty (prints nothing)
	// and the stdin case (prints counts without name).
	if len(output) > 0 {
		fmt.Println(output)
	}
}
