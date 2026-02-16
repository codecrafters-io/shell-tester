package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// MultipleCompletionsTestCase is a test case that:
// Sends text to the shell
// Asserts that the prompt line reflects the text
// Sends TAB
// Asserts that the expected reflection is printed to the screen (with a space after it)
// If any error occurs returns the error from the corresponding assertion
type MultipleCompletionsTestCase struct {
	// RawInput is the text to send to the shell
	RawInput string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// ExpectedAutocompletedReflectionHasNoSpace is true if
	// the expected reflection should have no space after it
	ExpectedAutocompletedReflectionHasNoSpace bool

	// CheckForBell is true if we should check for a bell
	CheckForBell bool

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// tabCount is the number of tabs to send after the input text
	TabCount int

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t MultipleCompletionsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Log the details of the text before sending it
	logTypedText(logger, t.RawInput)

	// Send the text to the shell
	if err := shell.SendText(t.RawInput); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	inputReflection := fmt.Sprintf("$ %s", t.RawInput)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: false,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", inputReflection)

	// // The space at the end of the reflection won't be present, so replace that assertion
	// asserter.PopAssertion()

	// Send TAB
	for i := range t.TabCount {
		shouldRingBell := i == 0 && t.CheckForBell
		logTab(logger, t.ExpectedReflection, shouldRingBell)

		// Node's readline doesn't register 2nd tab if sent instantly
		// Ref: CC-1689
		time.Sleep(5 * time.Millisecond)

		if err := shell.SendText("\t"); err != nil {
			return fmt.Errorf("Error sending text to shell: %v", err)
		}

		if shouldRingBell {
			// Assert no completions yet when the bell is received
			asserter.AddAssertion(assertions.EmptyLineAssertion{
				StayOnSameLine: true,
			})

			bellChannel := shell.VTBellChannel()
			asserter.AddAssertion(assertions.BellAssertion{
				BellChannel: bellChannel,
			})

			// Run the assertion, before sending the next tab
			if err := asserter.AssertWithoutPrompt(); err != nil {
				return err
			}
			logger.Successf("✓ Received bell")

			// Pop the bell assertion and empty line assertion after running
			asserter.PopAssertion()
			asserter.PopAssertion()
		}
	}

	inputTextReflection := t.ExpectedReflection
	// Space after autocomplete
	if !t.ExpectedAutocompletedReflectionHasNoSpace {
		inputTextReflection = fmt.Sprintf("%s ", t.ExpectedReflection)
	}

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputTextReflection,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", t.ExpectedReflection)

	// The space at the end of the reflection won't be present, so replace that assertion
	// asserter.PopAssertion()

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
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
