package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// AutocompleteTestCase is a test case that:
// Sends text to the shell
// Asserts that the prompt line reflects the typed text
// Sends TAB
// Asserts that the expected completion is printed to the screen (with a space after it)
// If any error occurs returns the error from the corresponding assertion
type AutocompleteTestCase struct {
	// PreviousInputOnLine is the string that is already present before RawInput is sent to the shell
	PreviousInputOnLine string

	// RawInput is the text to send to the shell
	RawInput string

	// ExpectedCompletion is the completion that is expected after the tab press
	ExpectedCompletion string

	// ExpectedCompletionHasNoSpace is true if
	// the expected completion should have no space after it
	ExpectedCompletionHasNoSpace bool

	// CheckForBell is true if we should check for a bell
	CheckForBell bool

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t AutocompleteTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Log the details of the text before sending it
	logTypedText(logger, t.RawInput)

	// Send the text to the shell
	if err := shell.SendText(t.RawInput); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	inputReflection := fmt.Sprintf("$ %s%s", t.PreviousInputOnLine, t.RawInput)
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

	// Remove the last assertion since the inputReflection won't be present in the next run
	// It will have been replaced by the expected completion
	asserter.PopAssertion()

	// Send TAB
	logTab(logger, t.ExpectedCompletion, false)
	if err := shell.SendText("\t"); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	expectedCompletion := fmt.Sprintf("$ %s", t.ExpectedCompletion)
	// Space after autocomplete
	if !t.ExpectedCompletionHasNoSpace {
		expectedCompletion += " "
	}
	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: expectedCompletion,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", t.ExpectedCompletion)

	// Remove the assertion after expected completion has been met
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

func logTab(logger *logger.Logger, expectedCompletion string, expectBell bool) {
	if expectBell {
		logger.Infof("Pressed %q (expecting bell to ring)", "<TAB>")
	} else {
		logger.Infof("Pressed %q (expecting autocomplete to %q)", "<TAB>", expectedCompletion)
	}
}

func logTypedText(logger *logger.Logger, text string) {
	logger.Infof("Typed %q", text)
}
