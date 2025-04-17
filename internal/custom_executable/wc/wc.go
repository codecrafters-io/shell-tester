package main

import (
	"bufio"
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
	lines, words, bytes int64 // Use int64 for potentially large files
}

func main() {
	// Use default flag set for simpler tools, or keep NewFlagSet if integrating into larger app
	// flag.Bool defines flags on the default command-line flag set
	lFlag := flag.Bool("l", false, "count lines")
	wFlag := flag.Bool("w", false, "count words")
	cFlag := flag.Bool("c", false, "count bytes") // Changed description

	flag.Parse() // Parse flags from os.Args[1:]

	// If no flags are provided, default to showing all counts
	showAll := !(*lFlag || *wFlag || *cFlag)
	if showAll {
		*lFlag = true
		*wFlag = true
		*cFlag = true
	}

	fileArgs := flag.Args()
	exitCode := 0 // Use 0 for success, 1 for failure

	if len(fileArgs) == 0 {
		// Read from stdin
		counts, err := countReader(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: error reading stdin: %v\n", err)
			os.Exit(1)
		}
		printCounts(counts, "-", *lFlag, *wFlag, *cFlag)
		os.Exit(0) // Explicitly exit cleanly
	}

	// Process each file
	total := counts{}
	filesProcessed := 0 // Track how many files were successfully processed for total line

	for _, filename := range fileArgs {
		f, err := os.Open(filename)
		if err != nil {
			// os.Open provides good errors for "No such file or directory" or permissions issues
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
			exitCode = 1 // Mark failure
			continue     // Skip to the next file
		}
		// Use defer for reliable closing
		defer f.Close() // Defer runs in LIFO order if multiple defers exist

		fileCounts, err := countReader(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: error reading %s: %v\n", filename, err)
			// Close might fail, but defer handles it. We should still report error & exit non-zero
			exitCode = 1
			// Optionally continue to next file or stop? wc usually continues.
			continue
		}

		// Use just the filename for display (as original code did)
		displayName := filepath.Base(filename)
		printCounts(fileCounts, displayName, *lFlag, *wFlag, *cFlag)

		total.lines += fileCounts.lines
		total.words += fileCounts.words
		total.bytes += fileCounts.bytes
		filesProcessed++

		// Close the file explicitly here *if not using defer*
		// err = f.Close()
		// if err != nil {
		//     fmt.Fprintf(os.Stderr, "wc: error closing %s: %v\n", filename, err)
		//     exitCode = 1
		// }
		// BUT: defer f.Close() placed after successful open is the idiomatic Go way
	}

	// Print total if multiple files were specified (even if some failed)
	// Standard wc prints total if more than one *argument* was given,
	// even if only one succeeded. This logic reflects that.
	if len(fileArgs) > 1 {
		printCounts(total, "total", *lFlag, *wFlag, *cFlag)
	}

	os.Exit(exitCode)
}

// countReader counts lines, words, and bytes from an io.Reader.
func countReader(r io.Reader) (counts, error) {
	var c counts
	counter := &byteCounter{} // Our byte counter

	// TeeReader sends reads to both the scanner's input and the byteCounter
	tee := io.TeeReader(r, counter)

	// Scanner will read from the tee, allowing byteCounter to count bytes
	scanner := bufio.NewScanner(tee)

	for scanner.Scan() {
		c.lines++
		line := scanner.Text() // Get text for word counting
		// strings.Fields splits by whitespace, handling multiple spaces correctly
		c.words += int64(len(strings.Fields(line)))
	}

	// After scanning finishes, check for scanning errors
	if err := scanner.Err(); err != nil {
		return counts{}, err // Return zero counts and the error
	}

	// The byte count comes from the byteCounter after the scanner consumed all input
	c.bytes = counter.count

	return c, nil
}

// printCounts formats and prints the counts for a given name (file or total).
func printCounts(c counts, name string, showLines, showWords, showBytes bool) {
	// Maintain consistent spacing
	prefix := ""
	if showLines {
		fmt.Printf("%s%8d", prefix, c.lines)
		prefix = " " // Add space before next field
	}
	if showWords {
		fmt.Printf("%s%8d", prefix, c.words)
		prefix = " "
	}
	if showBytes {
		fmt.Printf("%s%8d", prefix, c.bytes) // Changed from showChars
		prefix = " "
	}
	// Add space before the name if any counts were printed
	if prefix != "" {
		fmt.Printf(" %s\n", name)
	} else {
		// If no flags active (shouldn't happen with default logic, but safe)
		// Or if user explicitly uses e.g. `wc -l -w -c` on an empty file
		// We still print the name. If fields were printed, space is added above.
		// If no fields printed (e.g. `wc` with no flags and file is empty), print `0 0 0 filename`
		// If flags are off, print nothing but name? No, wc prints 0s if flags are specified.
		// Let's assume default logic ensures at least one flag is true if called.
		// If called with explicit flags off, this branch might not be hit.
		// The current logic handles the case where fields are printed correctly.
		// What if no counts are shown but name needs printing? Not possible with default flags.
		// Standard wc prints name even if counts are 0.
		fmt.Printf(" %s\n", name) // Ensure name is always printed after counts (or spaces for them)
	}
}
