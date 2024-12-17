package assertions

import (
	"fmt"
	"strings"

	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
)

// PromptAssertion verifies a prompt exists, and that there's no extra output after it.
type PromptAssertion struct {
	// ExpectedPrompt is the prompt expected to be displayed (example: "$ ")
	ExpectedPrompt string
}

func (t PromptAssertion) Inspect() string {
	return fmt.Sprintf("PromptAssertion (%q)", t.ExpectedPrompt)
}

func (t PromptAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	// We don't want to count the processed prompt as a complete row
	processedRowCount = 0

	rawRow := screenState[startRowIndex] // Could be nil?
	cleanedRow := virtual_terminal.BuildCleanedRow(rawRow)

	if !strings.EqualFold(cleanedRow, t.ExpectedPrompt) {
		return processedRowCount, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected prompt (%q) but received %q", t.ExpectedPrompt, cleanedRow),
		}
	}

	return processedRowCount, nil
}
