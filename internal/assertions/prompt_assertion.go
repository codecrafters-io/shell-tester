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

func (t PromptAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err error) {
	processedRowCount = 1

	if len(screenState) == 0 {
		return processedRowCount, fmt.Errorf("expected screen to have at least one row, but it was empty")
	}

	rawRow := screenState[startRowIndex]
	cleanedRow := utils.BuildCleanedRow(rawRow)

	if !strings.EqualFold(cleanedRow, t.expectedPrompt) {
		return processedRowCount, fmt.Errorf("expected prompt to be %q, but got %q", t.expectedPrompt, cleanedRow)
	}

	return processedRowCount, nil
}
