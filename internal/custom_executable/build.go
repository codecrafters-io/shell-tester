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

//go:embed main.txt
var content string

func ReplaceAndBuild(content, outputPath, placeholder, randomString string) error {
	// regex replace placeholder with random string
	content = strings.ReplaceAll(content, placeholder, randomString)

	// write content to file
	file, err := os.Create("tmp.go")
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: failed to create tmp.go: %w", err)
	}
	defer file.Close()
	file.WriteString(content)

	// Run go build command
	// ToDo: Remove log
	goCmdFullPath := path.Join(os.Getenv("TESTER_DIR"), "go")
	if goCmdFullPath == "go" {
		return fmt.Errorf("CodeCrafters Internal Error: Couldn't find packaged go command.\nTESTER_DIR: %s", os.Getenv("TESTER_DIR"))
	}
	buildCmd := exec.Command(goCmdFullPath, "build", "-o", outputPath, "tmp.go")
	buildCmd.Stdout = io.Discard
	buildCmd.Stderr = io.Discard
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: go build failed: %w", err)
	}

	err = os.Remove("tmp.go")
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: failed to remove tmp.go: %w", err)
	}

	return nil
}

func CreateExecutable(randomString, outputPath string) error {
	// Define the file path, output path, and placeholder
	placeholder := "PLACEHOLDER_RANDOM_STRING"

	// Call the replaceAndBuild function
	return ReplaceAndBuild(content, outputPath, placeholder, randomString)
}
