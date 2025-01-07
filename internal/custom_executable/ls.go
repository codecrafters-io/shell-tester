package custom_executable

import (
	"fmt"
	"os"
)

// Ls lists the contents of the specified directory(-ies)
// If no directory is provided, it lists the current directory
// Supports the -1 flag to list one file per line (default behavior)
func Ls(args []string) error {
	var dirArgs []string

	// Parse arguments
	for _, arg := range args {
		if arg == "-1" {
			continue // Skip the -1 flag since one-per-line is default
		}
		dirArgs = append(dirArgs, arg)
	}

	if len(dirArgs) == 0 {
		// ls defaults to current directory
		dirArgs = []string{"."}
	}

	for _, dir := range dirArgs {
		listOnePerLine(dir)
	}

	return nil
}

func listOnePerLine(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Print each file/directory name
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
