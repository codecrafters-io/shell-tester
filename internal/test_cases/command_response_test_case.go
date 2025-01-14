package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
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
	if err := shell.SendCommand(t.Command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.Command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   t.ExpectedOutput,
		FallbackPatterns: t.FallbackPatterns,
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
