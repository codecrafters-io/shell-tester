package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/screen_state"
)

// PromptAssertion verifies a prompt exists, and that there's no extra output after it.
type PromptAssertion struct {
	// ExpectedPrompt is the prompt expected to be displayed (example: "$ ")
	ExpectedPrompt string
}

func (t PromptAssertion) Inspect() string {
	return fmt.Sprintf("PromptAssertion (%q)", t.ExpectedPrompt)
}

func (t PromptAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	// We don't want to count the processed prompt as a complete row
	processedRowCount = 0

	rowString := screenState.GetRow(startRowIndex).String()

	if !strings.EqualFold(rowString, t.ExpectedPrompt) {
		return processedRowCount, &AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected prompt (%q) but received %q", t.ExpectedPrompt, rowString),
		}
	}

	return processedRowCount, nil
}
