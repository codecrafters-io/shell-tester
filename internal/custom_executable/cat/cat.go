package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Cat concatenates and prints files
// If no file is provided, it reads from stdin
func main() {
	flagSet := flag.NewFlagSet("cat", flag.ContinueOnError)
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		panic("cat: invalid option: " + err.Error())
	}

	fileArgs := flagSet.Args()

	if len(fileArgs) == 0 {
		// If no files provided, read from stdin
		if err := catFile(os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "cat: error reading stdin: %v\n", err)
		}
		return
	}

	// Process each file
	for _, file := range fileArgs {
		if !checkIfFileExists(file) {
			fmt.Fprintf(os.Stderr, "cat: %s: No such file or directory\n", file)
			continue
		}
		if err := processFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %s: %v\n", file, err)
			continue
		}
	}
}

func processFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return catFile(file)
}

func catFile(r io.Reader) error {
	_, err := io.Copy(os.Stdout, r)
	return err
}

func checkIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
