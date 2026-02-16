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

// CommandAutocompleteAndResponseTestCase is a test case that:
// Sends text to the shell
// Asserts that the prompt line reflects the text
// Sends TAB
// Asserts that the expected reflection is printed to the screen (with a space after it)
// Sends ENTER
// Asserts that the expected reflection is still present (with no space after it)
// If any error occurs returns the error from the corresponding assertion
type CommandAutocompleteAndResponseTestCase struct {
	// InputText is the text to send to the shell
	InputText string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// ExpectedAutocompletedReflectionHasNoSpace is true if
	// the expected reflection should have no space after it
	ExpectedAutocompletedReflectionHasNoSpace bool

	// Args is a list of arguments to pass to the command
	// Joins the args with spaces to form the expected output
	Args []string

	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandAutocompleteAndResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Log the details of the command before sending it
	logTypedText(logger, t.InputText)

	// Send the command to the shell
	if err := shell.SendText(t.InputText); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	inputReflection := fmt.Sprintf("$ %s", t.InputText)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", inputReflection)

	// The space at the end of the reflection won't be present, so replace that assertion
	asserter.PopAssertion()

	// Send TAB
	logTab(logger, t.ExpectedReflection, false)
	if err := shell.SendText("\t"); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflection)
	// Space after autocomplete
	if !t.ExpectedAutocompletedReflectionHasNoSpace {
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
	logger.Successf("✓ Prompt line matches %q", t.ExpectedReflection)

	// The space at the end of the reflection won't be present, so replace that assertion
	asserter.PopAssertion()

	// Send ENTER
	nextCommandToSend := "\n"
	if t.Args != nil {
		nextCommandToSend = strings.Join(t.Args, " ") + "\n"
	}
	logTypedText(logger, strings.TrimSpace(nextCommandToSend))
	logNewLine(logger)
	if err := shell.SendText(nextCommandToSend); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	// Assert the reflection again, after sending the enter key
	// This time, there won't be a space after the reflection
	commandReflection = fmt.Sprintf("$ %s %s", t.ExpectedReflection, strings.TrimSpace(nextCommandToSend))
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
	if t.ExpectedOutput != "" {
		logger.Successf("✓ Received %q", t.ExpectedOutput)
	}

	return nil
}
