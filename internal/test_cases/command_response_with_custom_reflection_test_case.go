package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandResponseWithCustomReflectionTestCase is a test case that:
// Sends a command to the shell
// Verifies that the ExpectedReflection is printed to the screen
// Reads the output from the shell, and verifies that it matches the expected output
// If any error occurs returns the error from the corresponding assertion
// NOTE: This test case uses SendCommandRaw to send the command to the shell
// So unless you pass `\t` or `\n` in the command, nothing interesting will happen
type CommandResponseWithCustomReflectionTestCase struct {
	// RawCommand is the command to send to the shell
	RawCommand string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandResponseWithCustomReflectionTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, skipCommandLogging bool) error {
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

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   t.ExpectedOutput,
		FallbackPatterns: t.FallbackPatterns,
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

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
