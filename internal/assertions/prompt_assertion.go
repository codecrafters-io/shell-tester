package assertions

import (
	"fmt"
	"strings"
)

// PromptTestCase verifies a prompt exists, and that there's no extra output after it.
type PromptAssertion struct {
	BaseAssertion

	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string
}

func NewPromptAssertion(screenAsserter *ScreenAsserter, rowIndex int, expectedPrompt string) PromptAssertion {
	return PromptAssertion{BaseAssertion: BaseAssertion{screenAsserter: screenAsserter, rowIndex: rowIndex}, expectedPrompt: expectedPrompt}
}

func (t PromptAssertion) Run() error {
	screen := t.screenAsserter.Shell.GetScreenState()
	if len(screen) == 0 {
		return fmt.Errorf("expected screen to have at least one row, but it was empty")
	}
	rawRow := screen[t.rowIndex]
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

func (t PromptAssertion) GetRowUpdateCount() int {
	return 0
}

func (t *PromptAssertion) UpdateRowIndex() {
	// Prompts are always on the same line, so we don't need to update the row index
	if t.ifUpdatedRowIndex {
		return
	}
	t.screenAsserter.UpdateRowIndex(t.GetRowUpdateCount())
	t.ifUpdatedRowIndex = true
	// fmt.Println("PromptAssertion.UpdateRowIndex() called, leading to row index", t.screenAsserter.GetRowIndex())
}
