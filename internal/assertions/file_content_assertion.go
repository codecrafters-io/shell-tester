package assertions

import (
	"fmt"
	"os"
)

// FileContentAssertion verifies a prompt exists, and that there's no extra output after it.
type FileContentAssertion struct {
	// FilePath is the path to the file to check
	FilePath string

	// ExpectedContent is the content expected to be in the file
	ExpectedContent string
}

func (t FileContentAssertion) Inspect() string {
	return fmt.Sprintf("FileContentAssertion (%q) with expected content (%q)", t.FilePath, t.ExpectedContent)
}

func (t FileContentAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	// We don't want to count the processed prompt as a complete row
	processedRowCount = 0

	fileContent, readErr := os.ReadFile(t.FilePath)
	if readErr != nil {
		return processedRowCount, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Failed to read file (%q): %v", t.FilePath, readErr),
		}
	}

	if string(fileContent) != t.ExpectedContent {
		return processedRowCount, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected (%q) contains (%q) but received %q", t.FilePath, t.ExpectedContent, fileContent),
		}
	}

	return processedRowCount, nil
}
