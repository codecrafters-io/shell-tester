package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type CommandWithMultilineResponseTestCase struct {
	// Command is the command to send to the shell
	Command string

	// MultiLineAssertion is the assertion to run
	MultiLineAssertion assertions.MultiLineAssertion

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipAssertPrompt is a flag to indicate that the prompt should not be asserted
	SkipPromptAssertion bool
}

func (t CommandWithMultilineResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if err := shell.SendCommand(t.Command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.Command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	asserter.AddAssertion(&t.MultiLineAssertion)

	if !t.SkipPromptAssertion {
		if err := asserter.AssertWithPrompt(); err != nil {
			return err
		}
	} else {
		if err := asserter.AssertWithoutPrompt(); err != nil {
			return err
		}
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
