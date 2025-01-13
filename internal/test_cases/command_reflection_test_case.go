package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandReflectionTestCase is a test case that:
// Sends a command to the shell
// Verifies that command is printed to the screen `$ <COMMAND>` (we expect the prompt to also be present)
// If any error occurs returns the error from the corresponding assertion
type CommandReflectionTestCase struct {
	// Command is the command to send to the shell
	Command string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandReflectionTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, skipSuccessMessage bool) error {
	return CommandWithCustomReflectionTestCase{
		RawCommand:          t.Command + "\n",
		ExpectedReflection:  t.Command,
		SuccessMessage:      t.SuccessMessage,
		SkipPromptAssertion: t.SkipPromptAssertion,
	}.Run(asserter, shell, logger, skipSuccessMessage, true)
}
