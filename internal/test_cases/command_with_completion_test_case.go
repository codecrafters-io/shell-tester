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

// CommandWithAttemptedCompletionTestCase is a test case that:
// Sends a partial command to the shell, sends TABS
// Asserts that the expected reflection is printed to the screen (with a space after it)
// Sends ENTER
// Asserts that the expected reflection is still present (with no space after it)
// If ExpectedOutput is provided,
// Reads the output from the shell, and verifies that it matches the expected output
// Else, it returns
// If any error occurs returns the error from the corresponding assertion
// NOTE: This test case uses SendCommandRaw to send the command to the shell
// So unless you pass `\t` or `\n` in the command, nothing interesting will happen
type CommandWithAttemptedCompletionTestCase struct {
	// RawCommand is the command to send to the shell
	RawCommand string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// ExpectedAutocompletedReflectionHasNoSpace is true if
	// the expected reflection should have no space after it
	ExpectedAutocompletedReflectionHasNoSpace bool

	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandWithAttemptedCompletionTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, skipCommandLogging bool) error {
	// TODO: Possibly panic if no tabs & newlines ?

	// Seperate the command into chars & tabs, newline
	hasEnterKey := t.RawCommand[len(t.RawCommand)-1] == '\n'
	hasTabKey := strings.Contains(t.RawCommand, "\t")

	rawCommand := t.RawCommand
	// Enter key is passed separately
	if hasEnterKey {
		rawCommand = t.RawCommand[:len(t.RawCommand)-1]
	}

	// Log the details of the command before sending it
	if !skipCommandLogging {
		LogCommandBeforeSending(logger, rawCommand, t.ExpectedReflection)
	}

	// Send the command to the shell
	if err := shell.SendCommandRaw(rawCommand); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflection)
	// Space after autocomplete
	if hasTabKey && !t.ExpectedAutocompletedReflectionHasNoSpace {
		commandReflection = fmt.Sprintf("$ %s ", t.ExpectedReflection)
	}
	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	if hasTabKey {
		logger.Successf("✓ Prompt line matches %q", t.ExpectedReflection)
	}
	// The space at the end of the reflection won't be present, so replace that assertion
	asserter.PopAssertion()

	if hasEnterKey {
		if err := shell.SendCommandRaw("\n"); err != nil {
			return fmt.Errorf("Error sending command to shell: %v", err)
		}
	}
	LogNewLine(logger)

	// Assert the reflection again, after sending the enter key
	// This time, there won't be a space after the reflection
	commandReflection = fmt.Sprintf("$ %s", t.ExpectedReflection)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	// If ExpectedOutput is provided, assert that the output matches
	if t.ExpectedOutput != "" {
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   t.ExpectedOutput,
			FallbackPatterns: t.FallbackPatterns,
		})
	}

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

func LogCommandBeforeSending(logger *logger.Logger, command string, expectedReflection string) {
	nonWhitespaceChars := ""
	tabs := ""
	for i := len(command) - 1; i >= 0; i-- {
		if command[i] == '\t' {
			tabs = string(command[i]) + tabs
		} else {
			nonWhitespaceChars = string(command[i]) + nonWhitespaceChars
		}
	}

	logger.Infof("Typed %q", nonWhitespaceChars)
	for range tabs {
		logger.Infof("Pressed %q", "<TAB>")
	}
}

func LogNewLine(logger *logger.Logger) {
	logger.Infof("Pressed %q", "<ENTER>")
}
