package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultLineCount = 10
	version          = "tail 0.1.0"
	helpText         = `Usage: tail [OPTION]... [FILE]...
Print the last 10 lines of each FILE to standard output.
With more than one FILE, precede each with a header giving the file name.
With no FILE, or when FILE is -, read standard input.

Mandatory arguments to long options are mandatory for short options too.
  -c, --bytes=K     output the last K bytes; or use -c +K to output
                    bytes starting with the Kth of each file
  -f, --follow[={name|descriptor}]
                    output appended data as the file grows
  -n, --lines=K     output the last K lines, instead of the last 10;
                    or use -n +K to output lines starting with the Kth
  -r                output lines in reverse order (unimplemented)
      --help        display this help and exit
      --version     output version information and exit`
)

// Represents the number and origin for line/byte count
type countParam struct {
	value         int64 // Using int64 for potential large file offsets
	fromBeginning bool  // True if count is relative to the start (e.g., +N)
}

type options struct {
	lineCount countParam
	byteCount countParam
	follow    bool
	reverse   bool
}

func parseOptions() (opts options, files []string, err error) {
	// Default options: last 10 lines
	opts = options{
		lineCount: countParam{value: defaultLineCount, fromBeginning: false},
		byteCount: countParam{value: -1, fromBeginning: false},
		follow:    false,
		reverse:   false,
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if arg == "--" { // End of options
				files = append(files, args[i+1:]...)
				break
			}

			switch {
			case arg == "--help":
				fmt.Println(helpText)
				os.Exit(0)
			case arg == "--version":
				fmt.Println(version)
				os.Exit(0)
			case arg == "-f", arg == "--follow":
				opts.follow = true
			case arg == "-r": // --reverse is not standard, only -r
				opts.reverse = true
			case strings.HasPrefix(arg, "-n") || strings.HasPrefix(arg, "--lines"):
				var valueStr string
				if arg == "-n" || arg == "--lines" {
					if i+1 >= len(args) {
						return opts, nil, fmt.Errorf("option requires an argument -- '%s'", arg)
					}
					i++
					valueStr = args[i]
				} else if strings.HasPrefix(arg, "-n=") {
					valueStr = arg[3:]
				} else if strings.HasPrefix(arg, "--lines=") {
					valueStr = arg[8:]
				} else { // Handle combined like -n5
					valueStr = arg[2:]
				}
				lineCount, err := parseCount(valueStr)
				if err != nil {
					return opts, nil, fmt.Errorf("illegal offset -- '%s'", valueStr)
				}
				opts.lineCount = lineCount
				opts.byteCount = countParam{value: -1} // Reset byte count if lines is specified later
			case strings.HasPrefix(arg, "-c") || strings.HasPrefix(arg, "--bytes"):
				var valueStr string
				if arg == "-c" || arg == "--bytes" {
					if i+1 >= len(args) {
						return opts, nil, fmt.Errorf("option requires an argument -- '%s'", arg)
					}
					i++
					valueStr = args[i]
				} else if strings.HasPrefix(arg, "-c=") {
					valueStr = arg[3:]
				} else if strings.HasPrefix(arg, "--bytes=") {
					valueStr = arg[8:]
				} else { // Handle combined like -c5
					valueStr = arg[2:]
				}
				byteCount, err := parseCount(valueStr)
				if err != nil {
					return opts, nil, fmt.Errorf("invalid number of bytes: '%s'", valueStr)
				}
				opts.byteCount = byteCount
				opts.lineCount = countParam{value: defaultLineCount, fromBeginning: false} // Reset line count
			default:
				// Handle combined single-character options like -n5, -c50, -f
				if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' {
					// Check for numeric suffix first (e.g., -5, -100)
					isNumericSuffix := false
					if numVal, err := strconv.ParseInt(arg[1:], 10, 64); err == nil {
						isNumericSuffix = true
						opts.lineCount = countParam{value: numVal, fromBeginning: false}
						opts.byteCount = countParam{value: -1}
						continue // Go to next argument
					}

					if !isNumericSuffix {
						// Iterate through combined flags like -fr
						for j := 1; j < len(arg); j++ {
							char := arg[j]
							switch char {
							case 'f':
								opts.follow = true
							case 'r':
								opts.reverse = true
							case 'c':
								// If -c is the last character, the next arg is the value
								if j == len(arg)-1 {
									if i+1 >= len(args) {
										return opts, nil, fmt.Errorf("option requires an argument -- 'c'")
									}
									i++
									valueStr := args[i]
									byteCount, err := parseCount(valueStr)
									if err != nil {
										return opts, nil, fmt.Errorf("invalid number of bytes: '%s'", valueStr)
									}
									opts.byteCount = byteCount
									opts.lineCount = countParam{value: defaultLineCount, fromBeginning: false} // Reset line count
								} else {
									// The value is immediately after -c (e.g., -c50)
									valueStr := arg[j+1:]
									byteCount, err := parseCount(valueStr)
									if err != nil {
										return opts, nil, fmt.Errorf("invalid number of bytes: '%s'", valueStr)
									}
									opts.byteCount = byteCount
									opts.lineCount = countParam{value: defaultLineCount, fromBeginning: false} // Reset line count
									j = len(arg)                                                               // Skip the rest of the combined flags
								}
							case 'n':
								// Similar logic for -n
								if j == len(arg)-1 {
									if i+1 >= len(args) {
										return opts, nil, fmt.Errorf("option requires an argument -- 'n'")
									}
									i++
									valueStr := args[i]
									lineCount, err := parseCount(valueStr)
									if err != nil {
										return opts, nil, fmt.Errorf("illegal offset -- '%s'", valueStr)
									}
									opts.lineCount = lineCount
									opts.byteCount = countParam{value: -1} // Reset byte count
								} else {
									valueStr := arg[j+1:]
									lineCount, err := parseCount(valueStr)
									if err != nil {
										return opts, nil, fmt.Errorf("illegal offset -- '%s'", valueStr)
									}
									opts.lineCount = lineCount
									opts.byteCount = countParam{value: -1} // Reset byte count
									j = len(arg)                           // Skip the rest
								}
							default:
								return opts, nil, fmt.Errorf("invalid option -- '%c'", char)
							}
						}
					}
				} else {
					return opts, nil, fmt.Errorf("invalid option -- '%s'", arg)
				}
			}
		} else {
			// If not an option flag or just "-", it's a file
			files = append(files, arg)
		}
	}

	// Final checks/adjustments
	if opts.byteCount.value >= 0 && opts.lineCount.value != defaultLineCount && !opts.lineCount.fromBeginning { // Byte count takes precedence if set explicitly
		// Only reset line count if it was the default, not explicitly set e.g. via -n+10 -c 5
		// Keep default line count if byte count is set
		// opts.lineCount = countParam{value: defaultLineCount, fromBeginning: false}
	} else if opts.byteCount.value < 0 && opts.lineCount.value == defaultLineCount { // If only default line count, ensure byte count is off
		// Keep the default line count active if byte count isn't set
		opts.byteCount.value = -1
	}

	// If -f is specified with stdin, POSIX says to ignore -f
	isStdinPresent := false
	if len(files) == 0 {
		isStdinPresent = true
	} else {
		for _, f := range files {
			if f == "-" {
				isStdinPresent = true
				break
			}
		}
	}
	if opts.follow && isStdinPresent {
		// Check if stdin is a pipe (not a FIFO)
		stat, _ := os.Stdin.Stat()
		if (stat.Mode()&os.ModeNamedPipe) != 0 && (stat.Mode()&os.ModeCharDevice) == 0 {
			// It's a pipe, ignore -f
			opts.follow = false
			// fmt.Fprintln(os.Stderr, "tail: warning: -f ignored when stdin is a pipe") // Optional warning
		}
	}

	return opts, files, nil
}

// parseCount parses a count value, handling optional leading '+'
func parseCount(value string) (countParam, error) {
	param := countParam{}
	if strings.HasPrefix(value, "+") {
		param.fromBeginning = true
		value = value[1:]
	}

	// GNU tail allows omitting the number for + (meaning +1)
	if param.fromBeginning && value == "" {
		param.value = 1
		return param, nil
	}

	count, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// Check for suffixes like k, M, G etc (GNU extension, not implementing yet)
		return param, err
	}

	if count < 0 {
		// This shouldn't happen with ParseInt unless the number is huge
		// but standard tail complains about negative numbers here
		return param, fmt.Errorf("value must be non-negative")
	}

	param.value = count
	return param, nil
}

func processFile(filename string, opts options, isFirst bool, multipleFilesPresent bool) error {
	var reader io.ReadSeeker // Use ReadSeeker for seeking (needed for -c/-n from end)
	var file *os.File
	var err error

	inputSourceIsStdin := false
	fileInfo, _ := os.Stdin.Stat()
	canSeekStdin := (fileInfo.Mode()&os.ModeCharDevice) == 0 && (fileInfo.Mode()&os.ModeNamedPipe) == 0

	if filename == "-" {
		if !canSeekStdin && (opts.byteCount.value > 0 || (opts.lineCount.value > 0 && !opts.lineCount.fromBeginning)) {
			// Cannot seek on stdin (like a pipe) for tailing from end
			// Need to read all input first for these modes.
			// Handle this special case within processBytes/processLines.
			// For now, we fall through, but the processing functions must handle it.
		}
		reader = os.Stdin
		filename = "standard input"
		inputSourceIsStdin = true
	} else {
		file, err = os.Open(filename)
		if err != nil {
			// Use a more standard error format
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("tail: %s: No such file or directory", filename)
			}
			return fmt.Errorf("tail: %v", err)
		}
		defer file.Close()
		reader = file
	}

	printHeader := !inputSourceIsStdin && multipleFilesPresent
	if printHeader && !isFirst {
		fmt.Println() // Add newline between outputs of different files
	}
	if printHeader {
		fmt.Printf("==> %s <==\n", filename)
	}

	if opts.reverse {
		// Reverse mode needs to read all lines first
		return processLinesReverse(reader, opts.lineCount)
	}

	if opts.byteCount.value >= 0 {
		err = processBytes(reader, opts.byteCount, inputSourceIsStdin)
	} else {
		err = processLines(reader, opts.lineCount, inputSourceIsStdin)
	}
	if err != nil {
		return err // Propagate errors from processing functions
	}

	// --- Follow Logic (-f) ---
	if opts.follow && !inputSourceIsStdin {
		return followFile(file)
	}

	return nil
}

// processBytes handles byte-based tailing
func processBytes(reader io.ReadSeeker, count countParam, isStdin bool) error {
	if count.value == 0 && count.fromBeginning { // -c +0 means print all
		count.value = 1 // Treat as -c +1
	}
	if count.value == 0 && !count.fromBeginning { // -c 0 means print nothing
		return nil
	}

	// Handle non-seekable input (like pipes) for counting from end
	canSeek := true
	if isStdin {
		info, _ := os.Stdin.Stat()
		canSeek = (info.Mode()&os.ModeCharDevice) == 0 && (info.Mode()&os.ModeNamedPipe) == 0
	}

	if !canSeek && !count.fromBeginning {
		// Read everything into memory to simulate tailing from end
		data, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("tail: error reading standard input: %v", err)
		}
		if int64(len(data)) > count.value {
			start := int64(len(data)) - count.value
			_, err = os.Stdout.Write(data[start:])
		} else {
			_, err = os.Stdout.Write(data)
		}
		return err
	}

	if count.fromBeginning {
		// Seek to the Kth byte (0-indexed, so K-1)
		_, err := reader.Seek(count.value-1, io.SeekStart)
		if err != nil {
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "invalid argument") { // EOF if seeking past end
				return nil // Nothing to print
			}
			return fmt.Errorf("tail: error seeking: %v", err)
		}
		// Copy the rest of the file to stdout
		_, err = io.Copy(os.Stdout, reader)
		return err
	} else {
		// Seek relative to the end
		fileSize, err := reader.Seek(0, io.SeekEnd)
		if err != nil {
			return fmt.Errorf("tail: error seeking to end: %v", err)
		}

		seekPos := int64(0)
		if fileSize > count.value {
			seekPos = fileSize - count.value
		}

		_, err = reader.Seek(seekPos, io.SeekStart)
		if err != nil {
			return fmt.Errorf("tail: error seeking to position %d: %v", seekPos, err)
		}

		// Copy from the calculated position to the end
		_, err = io.Copy(os.Stdout, reader)
		return err
	}
}

// processLines handles line-based tailing
func processLines(reader io.ReadSeeker, count countParam, isStdin bool) error {
	if count.value == 0 && count.fromBeginning { // -n +0 means print all
		count.value = 1 // Treat as -n +1
	}
	if count.value == 0 && !count.fromBeginning { // -n 0 means print nothing
		return nil
	}

	// Handle non-seekable input (like pipes)
	canSeek := true
	if isStdin {
		info, _ := os.Stdin.Stat()
		canSeek = (info.Mode()&os.ModeCharDevice) == 0 && (info.Mode()&os.ModeNamedPipe) == 0
	}

	if count.fromBeginning {
		// Read from start, skipping N-1 lines
		scanner := bufio.NewScanner(reader)
		lineNum := int64(1)
		for scanner.Scan() {
			if lineNum >= count.value {
				fmt.Println(scanner.Text())
			}
			lineNum++
		}
		return scanner.Err()
	} else {
		// Read last N lines
		if !canSeek {
			// Read all lines into memory for non-seekable input
			var lines []string
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("tail: error reading standard input: %v", err)
			}
			startIdx := 0
			if int64(len(lines)) > count.value {
				startIdx = len(lines) - int(count.value)
			}
			for _, line := range lines[startIdx:] {
				fmt.Println(line)
			}
			return nil
		} else {
			// Use seeking for regular files (more efficient for large files)
			// This is a simplified approach; a more robust one would read backwards in chunks.
			fileSize, err := reader.Seek(0, io.SeekEnd)
			if err != nil {
				return fmt.Errorf("tail: cannot seek: %v", err)
			}

			// Estimate starting point (very rough)
			// A better way involves reading backwards in chunks to find newlines
			seekPos := int64(0)
			if fileSize > count.value*80 { // Assuming avg line length 80
				seekPos = fileSize - count.value*80
			}
			_, err = reader.Seek(seekPos, io.SeekStart)
			if err != nil {
				return fmt.Errorf("tail: cannot seek: %v", err)
			}

			// Read lines from this point and keep the last N
			scanner := bufio.NewScanner(reader)
			var ringBuffer []string
			for scanner.Scan() {
				ringBuffer = append(ringBuffer, scanner.Text())
				if int64(len(ringBuffer)) > count.value {
					ringBuffer = ringBuffer[1:] // Remove oldest line
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("tail: error reading: %v", err)
			}

			// Print the lines in the buffer
			for _, line := range ringBuffer {
				fmt.Println(line)
			}
			return nil
		}
	}
}

// processLinesReverse handles -r option (unimplemented stub)
func processLinesReverse(reader io.Reader, count countParam) error {
	// Read all lines into memory
	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("tail: error reading input: %v", err)
	}

	startIdx := 0
	endIdx := len(lines)

	if count.value > 0 && !count.fromBeginning { // If -n N is given with -r, select last N before reversing
		if int64(len(lines)) > count.value {
			startIdx = len(lines) - int(count.value)
		}
	} else if count.value > 0 && count.fromBeginning {
		return fmt.Errorf("tail: cannot combine -n +K with -r") // Or ignore +K?
	}

	// Print lines in reverse order from the selected range
	for i := endIdx - 1; i >= startIdx; i-- {
		fmt.Println(lines[i])
	}

	return nil
}

// followFile handles the -f logic
func followFile(file *os.File) error {
	// For `follow`, we always start reading from the *current end* after the initial tail display.
	// So, regardless of where processLines/Bytes left off, seek to end now.
	endOffset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("tail: cannot seek to end of file %s: %v", file.Name(), err)
	}
	currentOffset := endOffset

	buffer := make([]byte, 4096) // Read in chunks

	for {
		// Check for truncation *before* reading
		stat, err := file.Stat()
		if err != nil {
			// Handle error: file might have been removed
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "tail: %s: file disappeared\n", file.Name())
				// Should we keep trying or exit? Standard tail often keeps trying.
				time.Sleep(1 * time.Second)
				continue // Keep trying
			} else {
				return fmt.Errorf("tail: error getting file stats for %s: %v", file.Name(), err)
			}
		}

		if stat.Size() < currentOffset {
			fmt.Fprintf(os.Stderr, "tail: %s: file truncated\n", file.Name())
			currentOffset, err = file.Seek(0, io.SeekEnd)
			if err != nil {
				return fmt.Errorf("tail: error seeking after truncation for %s: %v", file.Name(), err)
			}
		}

		// Try reading from the current offset
		// Ensure we're at the correct offset before reading
		_, err = file.Seek(currentOffset, io.SeekStart)
		if err != nil {
			return fmt.Errorf("tail: error seeking to offset %d for %s: %v", currentOffset, file.Name(), err)
		}

		n, err := file.Read(buffer)
		if n > 0 {
			_, writeErr := os.Stdout.Write(buffer[:n])
			if writeErr != nil {
				return fmt.Errorf("tail: error writing to stdout: %v", writeErr)
			}
			currentOffset += int64(n) // Update our position
		}

		if err != nil {
			if err == io.EOF {
				// End of file reached, wait for more data
				// Reduce sleep time for faster test execution
				time.Sleep(50 * time.Millisecond) // Poll interval reduced from 1s
				continue
			} else {
				// A real read error occurred
				return fmt.Errorf("tail: error reading %s: %v", file.Name(), err)
			}
		}
	}
	// This loop runs indefinitely until an error or interrupt
}

func main() {
	opts, files, err := parseOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tail: %v\n", err)
		os.Exit(1)
	}

	// If no files specified, use stdin
	if len(files) == 0 {
		files = append(files, "-")
	}

	exitCode := 0
	for i, filename := range files {
		err := processFile(filename, opts, i == 0, len(files) > 1)
		if err != nil {
			// Error message is already formatted by processFile or its callees
			fmt.Fprintf(os.Stderr, "%v\n", err)
			exitCode = 1
			// If -f is used, we might want to continue following other files?
			// Standard tail exits if any file causes an error initially.
			if !opts.follow {
				// No need to break if follow is active for other files
				// break // Exit on first error if not following?
			}
		}
	}

	// If following, the main loop would be elsewhere (or main would block)
	if !opts.follow || len(files) == 0 || (len(files) == 1 && files[0] == "-") {
		os.Exit(exitCode)
	} else {
		// If following files, keep the process alive (the actual following loop isn't here yet)
		select {} // Block forever
	}
}
