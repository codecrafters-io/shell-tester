package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type CommandWithMultilineResponseTestCase struct {
	// Command is the command to send to the shell
	Command string

	// ExpectedOutput is the expected output string to match against
	ExpectedOutput []string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SuccessMessage is the message to log in case of success
	SuccessMessage string
}

func (t CommandWithMultilineResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if err := shell.SendCommand(t.Command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.Command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	asserter.AddAssertion(assertions.MultiLineAssertion{
		ExpectedOutput:   t.ExpectedOutput,
		FallbackPatterns: t.FallbackPatterns,
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
