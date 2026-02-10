package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// AutocompleteTestCase is a test case that:
// Sends a prefix to the shell
// Asserts that the prompt line reflects the prefix
// Sends TAB
// Asserts that the expected reflection is printed to the screen (with/without a space as designated by ExpectedAutocompletedReflectionHasNoSpace)
// If any error occurs returns the error from the corresponding assertion
type AutocompleteTestCase struct {
	// TypedPrefix is the command to send to the shell
	TypedPrefix string

	// ExpectedPromptLineReflection is the custom reflection to use
	ExpectedPromptLineReflection string

	// ExpectedAutocompletedReflectionHasNoSpace is true if
	// the expected reflection should have no space after it
	ExpectedAutocompletedReflectionHasNoSpace bool

	// CheckForBell is true if we should check for a bell
	CheckForBell bool

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t AutocompleteTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Log the details of the command before sending it
	logCommand(logger, t.TypedPrefix)

	// Send the command to the shell
	if err := shell.SendCommandRaw(t.TypedPrefix); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	inputReflection := fmt.Sprintf("$ %s", t.TypedPrefix)
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
	logTab(logger, t.ExpectedPromptLineReflection, false)
	if err := shell.SendCommandRaw("\t"); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", t.ExpectedPromptLineReflection)
	// Space after autocomplete
	if !t.ExpectedAutocompletedReflectionHasNoSpace {
		commandReflection = fmt.Sprintf("$ %s ", t.ExpectedPromptLineReflection)
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
	logger.Successf("✓ Prompt line matches %q", t.ExpectedPromptLineReflection)
	// The space at the end of the reflection won't be present, so replace that assertion
	asserter.PopAssertion()

	if t.CheckForBell {
		bellChannel := shell.VTBellChannel()
		asserter.AddAssertion(assertions.BellAssertion{
			BellChannel: bellChannel,
		})
		// Run the assertion, before sending the enter key
		if err := asserter.AssertWithoutPrompt(); err != nil {
			return err
		}

		logger.Successf("✓ Received bell")
		// Pop the bell assertion after running
		asserter.PopAssertion()
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

	return nil
}

func logNewLine(logger *logger.Logger) {
	logger.Infof("Pressed %q", "<ENTER>")
}

func logTab(logger *logger.Logger, expectedReflection string, expectBell bool) {
	if expectBell {
		logger.Infof("Pressed %q (expecting bell to ring)", "<TAB>")
	} else {
		logger.Infof("Pressed %q (expecting autocomplete to %q)", "<TAB>", expectedReflection)
	}
}

func logCommand(logger *logger.Logger, command string) {
	logger.Infof("Typed %q", command)
}
