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

func (t *PromptAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err error) {
	if len(screenState) == 0 {
		return fmt.Errorf("expected screen to have at least one row, but it was empty")
	}

	rawRow := screenState[startRowIndex]
	cleanedRow := utils.BuildCleanedRow(rawRow)

	if !strings.EqualFold(cleanedRow, t.expectedPrompt) {
		return 0, fmt.Errorf("expected prompt to be %q, but got %q", t.expectedPrompt, cleanedRow)
	}

	return 0, nil
}

func (t *PromptAssertion) GetType() string {
	return "prompt"
}
