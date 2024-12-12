package logged_shell_asserter

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertion_collection"
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
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

func (a *LoggedShellAsserter) Assert() error {
	assertFn := func() error {
		return a.AssertionCollection.RunWithPromptAssertion(a.Shell.GetScreenState())
	}

	if readErr := a.Shell.ReadUntil(utils.AsBool(assertFn)); readErr != nil {
		if assertionErr := assertFn(); assertionErr != nil {
			a.logAssertionError(assertionErr)
			return fmt.Errorf("Assertion failed.")
		}
	}

	return nil
}

func (a *LoggedShellAsserter) onAssertionSuccess(startRowIndex int, processedRowCount int) {
	if processedRowCount == 0 {
		return
	}

	lastProcessedRowIndex := startRowIndex + processedRowCount - 1

	// fmt.Printf("debug: onAssertionSuccess called. startRowIndex: %d, processedRowCount: %d, lastProcessedRowIndex: %d\n", startRowIndex, processedRowCount, lastProcessedRowIndex)

	for i := 0; i < processedRowCount; i++ {
		// fmt.Printf("debug: logging1. i: %d, lastLoggedRowIndex: %d, processedRowCount: %d, lastProcessedRowIndex: %d", i, a.lastLoggedRowIndex, processedRowCount, lastProcessedRowIndex)
		row := a.Shell.GetScreenState()[a.lastLoggedRowIndex+i+1]
		// fmt.Printf("debug: logging2. i: %d, lastLoggedRowIndex: %d, processedRowCount: %d, lastProcessedRowIndex: %d, row: %q\n", i, a.lastLoggedRowIndex, processedRowCount, lastProcessedRowIndex, utils.BuildCleanedRow(row))
		a.Shell.LogOutput([]byte(utils.BuildCleanedRow(row)))
	}

	a.lastLoggedRowIndex = lastProcessedRowIndex

	if a.lastLoggedRowIndex > 3 {
		panic("end")
	}
}

func (a *LoggedShellAsserter) logAssertionError(err error) {
	// TODO: Log all shell output remaining
}
