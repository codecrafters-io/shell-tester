package assertions

import (
	"fmt"
	"strings"
)

// PromptTestCase verifies a prompt exists, and that there's no extra output after it.
type PromptAssertion struct {
	screenAsserter *ScreenAsserter

	rowIndex int

	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	// shouldOmitSuccessLog determines whether a success log should be written.
	// When re-using this test case within other higher-order test cases,
	// emitting success logs all the time can get pretty noisy.
	shouldOmitSuccessLog bool
}

func NewPromptAssertion(rowIndex int, expectedPrompt string) PromptAssertion {
	return PromptAssertion{rowIndex: rowIndex, expectedPrompt: expectedPrompt}
}

func NewSilentPromptAssertion(rowIndex int, expectedPrompt string) PromptAssertion {
	return PromptAssertion{rowIndex: rowIndex, expectedPrompt: expectedPrompt, shouldOmitSuccessLog: true}
}

func (t PromptAssertion) Run() error {
	rawRow := t.screenAsserter.Shell.GetScreenState()[t.rowIndex]
	cleanedRow := buildCleanedRow(rawRow)

	if !strings.EqualFold(cleanedRow, t.expectedPrompt) {
		return fmt.Errorf("expected prompt to be %q, but got %q", t.expectedPrompt, cleanedRow)
	}

	return nil
}

func (t PromptAssertion) WrappedRun() bool {
	// True if the prompt assertion is a success
	return t.Run() == nil
}
