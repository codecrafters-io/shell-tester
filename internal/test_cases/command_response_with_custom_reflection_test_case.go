package test_cases

import (
	"fmt"
	"regexp"
	"strings"

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
	hasEnterKey := t.RawCommand[len(t.RawCommand)-1] == '\n'
	hasTabKey := strings.Contains(t.RawCommand, "\t")
	rawCommand := t.RawCommand
	if hasEnterKey {
		rawCommand = t.RawCommand[:len(t.RawCommand)-1]
	}

	if !skipCommandLogging {
		LogCommandBeforeSending(logger, rawCommand, t.ExpectedReflection)
	}

	if err := shell.SendCommandRaw(rawCommand); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflection)
	// Space after autocomplete
	if hasTabKey {
		commandReflection = fmt.Sprintf("$ %s ", t.ExpectedReflection)
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
		StayOnSameLine: true,
	})

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	if hasTabKey {
		logger.Successf("✓ Prompt line matches %q", t.ExpectedReflection)
	}
	asserter.PopAssertion()

	if hasEnterKey {
		if err := shell.SendCommandRaw("\n"); err != nil {
			return fmt.Errorf("Error sending command to shell: %v", err)
		}
	}
	LogNewLine(logger)

	commandReflection = fmt.Sprintf("$ %s", t.ExpectedReflection)
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
