package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// PromptTestCase verifies a prompt exists, and that there's no extra output after it.
type PromptAssertion struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string
}

func NewPromptAssertion(expectedPrompt string) PromptAssertion {
	return PromptAssertion{expectedPrompt: expectedPrompt}
}

func (t PromptAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	// We don't want to count the processed prompt as a complete row
	processedRowCount = 0

	// TODO: Move these to assertion collection
	if len(screenState) == 0 {
		panic("CodeCrafters internal error: Expected screen state to have at least one row")
	}

	if startRowIndex >= len(screenState) {
		panic("CodeCrafters internal error: startRowIndex is larger than screenState rows")
	}

	rawRow := screenState[startRowIndex] // Could be nil?
	cleanedRow := utils.BuildCleanedRow(rawRow)

	if !strings.EqualFold(cleanedRow, t.expectedPrompt) {
		return processedRowCount, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected prompt (%q) but received %q", t.expectedPrompt, cleanedRow),
		}
	}

	return processedRowCount, nil
}
