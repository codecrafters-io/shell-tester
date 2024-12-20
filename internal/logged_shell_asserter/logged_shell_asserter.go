package logged_shell_asserter

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertion_collection"
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
)

type LoggedShellAsserter struct {
	Shell               *shell_executable.ShellExecutable
	AssertionCollection *assertion_collection.AssertionCollection

	lastLoggedRowIndex int
}

func NewLoggedShellAsserter(shell *shell_executable.ShellExecutable) *LoggedShellAsserter {
	assertionCollection := assertion_collection.NewAssertionCollection()

	asserter := &LoggedShellAsserter{
		Shell:               shell,
		AssertionCollection: assertionCollection,
		lastLoggedRowIndex:  -1,
	}

	assertionCollection.OnAssertionSuccess = asserter.onAssertionSuccess

	return asserter
}

func (a *LoggedShellAsserter) AddAssertion(assertion assertions.Assertion) {
	a.AssertionCollection.AddAssertion(assertion)
}

func (a *LoggedShellAsserter) AssertWithPrompt() error {
	return a.assert(false)
}

func (a *LoggedShellAsserter) AssertWithoutPrompt() error {
	return a.assert(true)
}

func (a *LoggedShellAsserter) assert(withoutPrompt bool) error {
	var assertFn func() *assertions.AssertionError

	if withoutPrompt {
		assertFn = func() *assertions.AssertionError {
			return a.AssertionCollection.RunWithoutPromptAssertion(a.Shell.GetScreenState())
		}
	} else {
		assertFn = func() *assertions.AssertionError {
			return a.AssertionCollection.RunWithPromptAssertion(a.Shell.GetScreenState())
		}
	}

	conditionFn := func() bool {
		return assertFn() == nil
	}

	if readErr := a.Shell.ReadUntil(conditionFn); readErr != nil {
		if assertionErr := assertFn(); assertionErr != nil {
			a.logAssertionError(*assertionErr)
			return fmt.Errorf("Assertion failed.")
		}
	}

	return nil
}

func (a *LoggedShellAsserter) onAssertionSuccess(startRowIndex int, processedRowCount int) {
	if processedRowCount == 0 || startRowIndex <= a.lastLoggedRowIndex {
		return
	}

	for i := 0; i < processedRowCount; i++ {
		row := a.Shell.GetScreenState()[a.lastLoggedRowIndex+i+1]
		a.Shell.LogOutput([]byte(virtual_terminal.BuildCleanedRow(row)))
	}

	a.lastLoggedRowIndex += processedRowCount
}

func (a *LoggedShellAsserter) logAssertionError(err assertions.AssertionError) {
	a.logRows(a.lastLoggedRowIndex+1, err.ErrorRowIndex+1)
	a.Shell.GetLogger().Errorf("%s", err.Message)
	a.logRows(err.ErrorRowIndex+1, len(a.Shell.GetScreenState()))
}

func (a *LoggedShellAsserter) LogRemainingOutput() {
	startRowIndex := a.lastLoggedRowIndex + 1
	endRowIndex := len(a.Shell.GetScreenState())
	a.logRows(startRowIndex, endRowIndex)
}

func (a *LoggedShellAsserter) logRows(startRowIndex int, endRowIndex int) {
	for i := startRowIndex; i < endRowIndex; i++ {
		rawRow := a.Shell.GetScreenState()[i]
		cleanedRow := virtual_terminal.BuildCleanedRow(rawRow)
		if len(cleanedRow) > 0 {
			a.Shell.LogOutput([]byte(cleanedRow))
		}
	}
}

func (a *LoggedShellAsserter) GetLastLoggedRowIndex() int {
	return a.lastLoggedRowIndex
}
