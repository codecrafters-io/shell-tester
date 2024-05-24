package custom_executable

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
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
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	file.WriteString(content)

	// Run go build command
	buildCmd := exec.Command("go", "build", "-o", outputPath, "tmp.go")
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to run go build: %w", err)
	}

	err = os.Remove("tmp.go")
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

func CreateExecutable(randomString, outputPath string) error {
	// Define the file path, output path, and placeholder
	placeholder := "PLACEHOLDER_RANDOM_STRING"

	// Call the replaceAndBuild function
	return ReplaceAndBuild(content, outputPath, placeholder, randomString)
}
