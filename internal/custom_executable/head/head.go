package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	defaultLineCount = 10
	version          = "head 0.1.0"
	helpText         = `Usage: head [OPTION]... [FILE]...
Print the first 10 lines of each FILE to standard output.
With more than one FILE, precede each with a header giving the file name.
With no FILE, or when FILE is -, read standard input.

Mandatory arguments to long options are mandatory for short options too.
  -c, --bytes=K     print the first K bytes of each file
  -n, --lines=K     print the first K lines instead of the first 10
  --help               display this help and exit
  --version            output version information and exit`
)

type options struct {
	lineCount   int
	byteCount   int
	quiet       bool
	verbose     bool
	invertCount bool
}

func parseOptions() (opts options, files []string, err error) {
	// Default options
	opts = options{
		lineCount: defaultLineCount,
		byteCount: -1,
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check if it's an option flag
		if strings.HasPrefix(arg, "-") {
			switch {
			case arg == "--help":
				fmt.Println(helpText)
				os.Exit(0)
			case arg == "--version":
				fmt.Println(version)
				os.Exit(0)
			case strings.HasPrefix(arg, "-n=") || strings.HasPrefix(arg, "--lines="):
				var value string
				if strings.HasPrefix(arg, "-n=") {
					value = arg[3:]
				} else {
					value = arg[8:]
				}
				lineCount, err := parseCount(value)
				if err != nil {
					return opts, nil, fmt.Errorf("invalid line count: %v", err)
				}
				opts.lineCount = lineCount
			case strings.HasPrefix(arg, "-c=") || strings.HasPrefix(arg, "--bytes="):
				var value string
				if strings.HasPrefix(arg, "-c=") {
					value = arg[3:]
				} else {
					value = arg[8:]
				}
				byteCount, err := parseCount(value)
				if err != nil {
					return opts, nil, fmt.Errorf("invalid byte count: %v", err)
				}
				opts.byteCount = byteCount
			case arg == "-n" || arg == "--lines":
				if i+1 >= len(args) {
					return opts, nil, fmt.Errorf("option requires an argument -- '%s'", arg)
				}
				i++
				lineCount, err := parseCount(args[i])
				if err != nil {
					return opts, nil, fmt.Errorf("invalid line count: %v\n", err)
				}
				if lineCount < 0 {
					return opts, nil, fmt.Errorf("illegal line count -- %d\n", lineCount)
				}
				opts.lineCount = lineCount
			case arg == "-c" || arg == "--bytes":
				if i+1 >= len(args) {
					return opts, nil, fmt.Errorf("option requires an argument -- '%s'", arg)
				}
				i++
				byteCount, err := parseCount(args[i])
				if err != nil {
					return opts, nil, fmt.Errorf("invalid byte count: %v", err)
				}
				if byteCount < 0 {
					return opts, nil, fmt.Errorf("illegal byte count -- %d", byteCount)
				}
				opts.byteCount = byteCount
			default:
				// Handle single-character combined options like -n5 or -qv
				if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' {
					for j := 1; j < len(arg); j++ {
						switch arg[j] {
						case 'n':
							// If -n is the last character, the next arg is the value
							if j == len(arg)-1 {
								if i+1 >= len(args) {
									return opts, nil, fmt.Errorf("option requires an argument -- 'n'")
								}
								i++
								lineCount, err := parseCount(args[i])
								if err != nil {
									return opts, nil, fmt.Errorf("invalid line count: %v", err)
								}
								opts.lineCount = lineCount
							} else {
								// The value is immediately after -n
								lineCount, err := parseCount(arg[j+1:])
								if err != nil {
									return opts, nil, fmt.Errorf("illegal line count: %v", err)
								}
								opts.lineCount = lineCount
								j = len(arg) // Skip the rest
							}
						case 'c':
							// If -c is the last character, the next arg is the value
							if j == len(arg)-1 {
								if i+1 >= len(args) {
									return opts, nil, fmt.Errorf("option requires an argument -- 'c'")
								}
								i++
								byteCount, err := parseCount(args[i])
								if err != nil {
									return opts, nil, fmt.Errorf("invalid byte count: %v", err)
								}
								if byteCount < 0 {
									return opts, nil, fmt.Errorf("illegal byte count -- %d", byteCount)
								}
								opts.byteCount = byteCount
							} else {
								// The value is immediately after -c
								byteCount, err := parseCount(arg[j+1:])
								if err != nil {
									return opts, nil, fmt.Errorf("invalid byte count: %v", err)
								}
								opts.byteCount = byteCount
								j = len(arg) // Skip the rest
							}
						default:
							return opts, nil, fmt.Errorf("invalid option -- '%c'", arg[j])
						}
					}
				} else {
					return opts, nil, fmt.Errorf("invalid option -- '%s'", arg)
				}
			}
		} else {
			// If not an option flag, it's a file
			files = append(files, arg)
		}
	}

	return opts, files, nil
}

// parseCount parses a count value which may have a leading '-' for inversion
// or a suffix like k, M, etc.
func parseCount(value string) (int, error) {
	invertCount := false
	if strings.HasPrefix(value, "-") {
		invertCount = true
		value = value[1:]
	}

	// Check for multiplier suffixes
	multiplier := 1
	if len(value) > 1 {
		lastChar := value[len(value)-1]
		switch lastChar {
		case 'b':
			value = value[:len(value)-1]
			multiplier = 512
		case 'k':
			value = value[:len(value)-1]
			multiplier = 1000
		case 'K':
			value = value[:len(value)-1]
			multiplier = 1024
		case 'm':
			value = value[:len(value)-1]
			multiplier = 1000 * 1000
		case 'M':
			value = value[:len(value)-1]
			multiplier = 1024 * 1024
		case 'g':
			value = value[:len(value)-1]
			multiplier = 1000 * 1000 * 1000
		case 'G':
			value = value[:len(value)-1]
			multiplier = 1024 * 1024 * 1024
		}
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	count *= multiplier
	if invertCount {
		count = -count
	}

	return count, nil
}

func processFile(filename string, opts options, isFirst bool, multipleFilesPresent bool) error {
	var reader io.Reader
	var file *os.File
	var err error

	// Determine input source (stdin or file)
	inputSourceIsStdin := false

	if filename == "-" {
		reader = os.Stdin
		filename = "standard input"
		inputSourceIsStdin = true
	} else {
		file, err = os.Open(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("head: %s: No such file or directory\n", filename)
			}
			return fmt.Errorf("head: %s: %v\n", filename, err)
		}
		defer file.Close()
		reader = file
	}

	// Print header if needed
	printHeader := !inputSourceIsStdin && (multipleFilesPresent)
	if printHeader && !isFirst {
		fmt.Println()
	}
	if printHeader {
		fmt.Printf("==> %s <==\n", filename)
	}

	// Process by bytes if -c is specified
	if opts.byteCount >= 0 {
		return processBytes(reader, opts.byteCount)
	}

	// Otherwise process by lines
	return processLines(reader, opts.lineCount)
}

func processBytes(reader io.Reader, byteCount int) error {
	if byteCount == 0 {
		return nil
	}

	if byteCount > 0 {
		// Print first N bytes
		buffer := make([]byte, byteCount)
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		os.Stdout.Write(buffer[:n])
		return nil
	}

	// For negative byte count (all but last N bytes)
	// We need to read the entire file into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Print all but the last N bytes
	absCount := -byteCount
	if len(data) > absCount {
		os.Stdout.Write(data[:len(data)-absCount])
	}
	return nil
}

func processLines(reader io.Reader, lineCount int) error {
	if lineCount == 0 {
		return nil
	}

	scanner := bufio.NewScanner(reader)

	if lineCount > 0 {
		// Print first N lines
		linesRead := 0
		for scanner.Scan() && linesRead < lineCount {
			fmt.Println(scanner.Text())
			linesRead++
		}
		return scanner.Err()
	}

	// For negative line count (all but last N lines)
	// We need to keep a rolling buffer of the last N lines
	absCount := -lineCount
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > absCount {
			fmt.Println(lines[0])
			lines = lines[1:]
		}
	}
	return scanner.Err()
}

func main() {
	opts, files, err := parseOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "head: %v", err)
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
			fmt.Fprintf(os.Stderr, "%v", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
