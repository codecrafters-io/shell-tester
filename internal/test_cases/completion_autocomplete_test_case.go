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
	logTabForCompletion(logger, t.ExpectedCompletion, false)
	if err := shell.SendText("\t"); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	expectedCompletion := fmt.Sprintf("$ %s", t.ExpectedCompletion)

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: expectedCompletion,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// If the completion does not change the prompt line: notify that prompt line is unchanged
	if t.ExpectedCompletion != t.PreviousInputOnLine+t.RawInput {
		logger.Successf("✓ Prompt line matches %q", t.ExpectedCompletion)
	} else {
		logger.Successf("✓ Prompt line unchanged after <TAB> press")
	}

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

func logTabForCompletion(logger *logger.Logger, expectedCompletion string, expectBell bool) {
	if expectBell {
		logger.Infof("Pressed %q (expecting bell to ring)", "<TAB>")
		return
	}

	if expectedCompletion[len(expectedCompletion)-1] == ' ' {
		expectedCompletionWithoutSpace := expectedCompletion[:len(expectedCompletion)-1]
		logger.Infof("Pressed %q (expecting autocomplete to %q followed by a space)", "<TAB>", expectedCompletionWithoutSpace)
		return
	}

	logger.Infof("Pressed %q (expecting autocomplete to %q)", "<TAB>", expectedCompletion)
}

func logTypedText(logger *logger.Logger, text string) {
	if len(text) == 0 {
		return
	}

	hasEndingSpace := text[len(text)-1] == ' '
	hasStartingSpace := text[0] == ' '

	if (len(text) == 1) || (!hasEndingSpace && !hasStartingSpace) {
		logger.Infof("Typed %q", text)
		return
	}

	if hasEndingSpace && hasStartingSpace {
		logger.Infof("Typed <SPACE>, followed by %q, followed by <SPACE>", text[1:len(text)-1])
		return
	}

	if hasEndingSpace {
		logger.Infof("Typed %q followed by a <SPACE>", text[:len(text)-1])
		return
	}

	if hasStartingSpace {
		logger.Infof("Typed <SPACE>, followed by %q", text[1:])
		return
	}
}
