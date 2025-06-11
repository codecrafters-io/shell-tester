package logged_shell_asserter

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertion_collection"
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
)

// INITIAL_READ_TIMEOUT is used for the first prompt read, where we want
// to be more lenient, and allow user's shells to start up properly
const INITIAL_READ_TIMEOUT = 5000 * time.Millisecond
const SUBSEQUENT_READ_TIMEOUT = 2000 * time.Millisecond

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

func (a *LoggedShellAsserter) StartShellAndAssertPrompt(skipSuccessMessage bool) error {
	if err := a.Shell.Start(); err != nil {
		return err
	}

	if err := a.AssertWithPromptAndLongerTimeout(); err != nil {
		return err
	}

	if !skipSuccessMessage {
		a.Shell.GetLogger().Successf("✓ Received prompt ($ )")
	}

	// .NET ReadLine() method seems to have a bug where it prints the command twice
	// in certain cases. This sleep is a workaround for that. Refer to: CC-1576
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (a *LoggedShellAsserter) AddAssertion(assertion assertions.Assertion) {
	a.AssertionCollection.AddAssertion(assertion)
}

func (a *LoggedShellAsserter) PopAssertion() assertions.Assertion {
	return a.AssertionCollection.PopAssertion()
}

func (a *LoggedShellAsserter) AssertWithPrompt() error {
	return a.assert(false, SUBSEQUENT_READ_TIMEOUT)
}

func (a *LoggedShellAsserter) AssertWithoutPrompt() error {
	return a.assert(true, SUBSEQUENT_READ_TIMEOUT)
}

func (a *LoggedShellAsserter) AssertWithPromptAndLongerTimeout() error {
	return a.assert(false, INITIAL_READ_TIMEOUT)
}

func (a *LoggedShellAsserter) assert(withoutPrompt bool, readTimeout time.Duration) error {
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

	if readErr := a.Shell.ReadUntilConditionOrTimeout(conditionFn, readTimeout); readErr != nil {
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

	a.logRowsUntil(a.lastLoggedRowIndex + processedRowCount)
}

func (a *LoggedShellAsserter) logAssertionError(err assertions.AssertionError) {
	a.logRowsUntil(err.ErrorRowIndex)
	a.Shell.GetLogger().Errorf("%s", err.Message)
	a.LogRemainingOutput()
}

func (a *LoggedShellAsserter) LogRemainingOutput() {
	screenState := a.Shell.GetScreenState()
	cursorRowIndex, cursorColIndex := a.Shell.GetCursorPosition()
	lastRowIndex := len(screenState) - 1
	lastRowContents := virtual_terminal.BuildCleanedRow(screenState[lastRowIndex])

	// If the last row is (a) empty and (b) has the cursor at the start we can skip logging it
	if lastRowContents == "" && cursorRowIndex == lastRowIndex && cursorColIndex == 0 {
		a.logRowsUntil(lastRowIndex - 1)
	} else {
		a.logRowsUntil(lastRowIndex)
	}
}

func (a *LoggedShellAsserter) logRowsUntil(endRowIndex int) {
	for i := a.lastLoggedRowIndex + 1; i <= endRowIndex; i++ {
		// If we've reached the end of the screen state we can stop logging
		if i > len(a.Shell.GetScreenState())-1 {
			return
		}

		rawRow := a.Shell.GetScreenState()[i]
		cleanedRow := virtual_terminal.BuildCleanedRow(rawRow)
		a.Shell.LogOutput([]byte(cleanedRow))
		a.lastLoggedRowIndex = i
	}
}

func (a *LoggedShellAsserter) GetLastLoggedRowIndex() int {
	return a.lastLoggedRowIndex
}
