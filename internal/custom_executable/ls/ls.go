package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

// Ls lists the contents of the specified directory(-ies)
// If no directory is provided, it lists the current directory
// Supports the -1 flag to list one file per line (default behavior)
func main() {
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	// ls -1 is the default behavior
	_ = flagSet.Bool("1", false, "list one file per line")
	// Parse flags, would return error if any other flags are provided
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		panic("CodeCrafters Internal Error: " + err.Error())
	}

	dirArgs := flagSet.Args()

	if len(dirArgs) == 0 {
		// ls defaults to current directory
		dirArgs = []string{"."}
	}

	// If multiple directories are provided, ls sorts them
	sort.Strings(dirArgs)
	dirsWhichExist := []string{}
	for _, dir := range dirArgs {
		if !checkIfDirectoryExists(dir) {
			fmt.Printf("ls: %s: No such file or directory\n", dir)
			continue
		}
		dirsWhichExist = append(dirsWhichExist, dir)
	}

	multipleDirsPresent := len(dirsWhichExist) > 1
	for i, dir := range dirsWhichExist {
		if multipleDirsPresent {
			fmt.Printf("%s:\n", dir)
		}
		listOnePerLine(dir)
		// New line between each directory's entries
		if multipleDirsPresent && i < len(dirsWhichExist)-1 {
			fmt.Println()
		}
	}
}

func checkIfDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func listOnePerLine(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("ls: %s: No such file or directory\n", path)
		return
	}

	// Print each file/directory name
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
