package test_cases

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandResponseTestCase is a test case that:
// Sends a command to the shell
// Verifies that command is printed to the screen `$ <COMMAND>` (we expect the prompt to also be present)
// Reads the output from the shell, and verifies that it matches the expected output
// If any error occurs returns the error from the corresponding assertion
type CommandResponseTestCase struct {
	// Command is the command to send to the shell
	Command string

	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SuccessMessage is the message to log in case of success
	SuccessMessage string
}

func (t CommandResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	return CommandResponseWithCustomReflectionTestCase{
		RawCommand:         t.Command + "\n",
		ExpectedReflection: t.Command,
		ExpectedOutput:     t.ExpectedOutput,
		FallbackPatterns:   t.FallbackPatterns,
		SuccessMessage:     t.SuccessMessage,
	}.Run(asserter, shell, logger, true)
}
