package logged_shell_asserter

import (
	"github.com/codecrafters-io/shell-tester/internal/assertion_collection"
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

type LoggedShellAsserter struct {
	Shell               *shell_executable.ShellExecutable
	AssertionCollection *assertion_collection.AssertionCollection
}

func NewLoggedShellAsserter(shell *shell_executable.ShellExecutable) *LoggedShellAsserter {
	return &LoggedShellAsserter{
		Shell:               shell,
		AssertionCollection: assertion_collection.NewAssertionCollection(),
	}
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
			return assertionErr
		}
	}

	return nil
}
