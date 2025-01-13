package test_cases

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandWithCustomReflectionTestCase is a test case that:
// Sends a command to the shell
// Verifies that ExpectedReflection is printed to the screen
// If any error occurs returns the error from the corresponding assertion
// NOTE: This test case uses SendCommandRaw to send the command to the shell
// So unless you pass `\t` or `\n` in the command, nothing interesting will happen
type CommandWithCustomReflectionTestCase struct {
	// RawCommand is the command to send to the shell
	RawCommand string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandWithCustomReflectionTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, skipSuccessMessage bool, skipCommandLogging bool) error {
	if !skipCommandLogging {
		LogCommandBeforeSending(logger, t.RawCommand)
	}
	if err := shell.SendCommandRaw(t.RawCommand); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflection)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	var assertFuncToRun func() error
	if t.SkipPromptAssertion {
		assertFuncToRun = asserter.AssertWithoutPrompt
	} else {
		assertFuncToRun = asserter.AssertWithPrompt
	}

	if err := assertFuncToRun(); err != nil {
		return err
	}

	if !skipSuccessMessage {
		logger.Successf("%s", t.SuccessMessage)
	}
	return nil
}

func LogCommandBeforeSending(logger *logger.Logger, command string) {
	command = strings.ReplaceAll(command, " ", "·")
	command = strings.ReplaceAll(command, "\t", "⇥ ") // →
	command = strings.ReplaceAll(command, "\n", "⏎")
	logger.Infof("[setup] %s", command)
}
