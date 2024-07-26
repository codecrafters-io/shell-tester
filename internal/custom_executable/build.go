package custom_executable

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

//go:embed c.txt
var content string

func ReplaceAndBuild(content, outputPath, placeholder, randomString string) error {
	// regex replace placeholder with random string
	content = strings.ReplaceAll(content, placeholder, randomString)

	// write content to file
	file, err := os.Create("tmp.c")
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: failed to create tmp.c: %w", err)
	}
	defer file.Close()
	file.WriteString(content)
	defer func() {
		// Remove the file even if the build fails
		if err := os.Remove("tmp.c"); err != nil {
			fmt.Printf("CodeCrafters Internal Error: failed to remove tmp.c: %v\n", err)
		}
	}()

	// Run tcc build command
	// tcc is included in the tester directory
	tccCmdFullPath := path.Join(os.Getenv("TESTER_DIR"), "tcc")
	if tccCmdFullPath == "" {
		return fmt.Errorf("CodeCrafters Internal Error: Couldn't find tcc command")
	}
	buildCmd := exec.Command(tccCmdFullPath, "tmp.c", "-o", outputPath)
	buildCmd.Stdout = io.Discard
	buildCmd.Stderr = io.Discard
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: tcc build failed: %w", err)
	}

	return nil
}

func CreateExecutable(randomString, outputPath string) error {
	// Define the file path, output path, and placeholder
	placeholder := "PLACEHOLDER_RANDOM_STRING"

	// Call the replaceAndBuild function
	return ReplaceAndBuild(content, outputPath, placeholder, randomString)
}
