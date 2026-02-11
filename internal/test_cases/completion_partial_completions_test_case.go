package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// PartialCompletionsTestCase is a test case that:
// Sends a command to the shell
// Asserts that the prompt line reflects the command
// for each partial auto-completion:
// Sends TAB
// Asserts that the expected reflection is printed to the screen
// And sends the subsequent input if it is non-empty
// If any error occurs returns the error from the corresponding assertion
type PartialCompletionsTestCase struct {
	// TODO refactor: combine input and reflections in one structure
	// This needs this test case's usage to be changed across the previous extension as well
	// skipping the refactor for now

	// Inputs is the list of inputs to send to the shell
	// They are send one by one, interleaved with TABs
	// The shell is expected to auto-complete expected reflections
	// If any of the inputs is an empty string, its assertion is skipped
	Inputs []string

	// ExpectedReflections is the list of expected reflections to use
	ExpectedReflections []string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool

	// ExpectedLastReflectionHasNoSpace should be true if no space is expected
	// after the last reflection appears
	ExpectedLastReflectionHasNoSpace bool
}

func (t PartialCompletionsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if len(t.Inputs) != len(t.ExpectedReflections) {
		panic("Inputs and ExpectedReflections must have the same length")
	}

	// The entire flow is repeated for each input & expected reflection
	for idx := 0; idx < len(t.ExpectedReflections); idx++ {
		if err := t.runInputReflectionForIdx(asserter, shell, logger, idx); err != nil {
			return err
		}

		if err := t.runTabCompletionForIdx(asserter, shell, logger, idx); err != nil {
			return err
		}
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

// runInputReflection runs the input single line assertion for the input at index 'idx'
func (t PartialCompletionsTestCase) runInputReflectionForIdx(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, idx int) error {
	// We don't need to run the test case if the input is empty
	if t.Inputs[idx] == "" {
		return nil
	}

	// Log the details of the command before sending it
	logTypedText(logger, t.Inputs[idx])

	// Send the command to the shell
	if err := shell.SendTextRaw(t.Inputs[idx]); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	// The prompt line will not just show the subsequent input,
	// but the previous reflection concatenated with the current input, if any
	prevInput := ""
	if idx > 0 {
		prevInput = t.ExpectedReflections[idx-1]
	}

	inputReflection := fmt.Sprintf("$ %s", prevInput+t.Inputs[idx])
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the tab key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	asserter.PopAssertion()
	logger.Successf("✓ Prompt line matches %q", inputReflection)
	return nil
}

func (t PartialCompletionsTestCase) runTabCompletionForIdx(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, idx int) error {
	// Send TAB
	logTab(logger, t.ExpectedReflections[idx], false)
	if err := shell.SendTextRaw("\t"); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	// For all partial auto-completions, we expect *NO* space at the end
	commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflections[idx])

	// For the last auto-completion, we expect a space at the end if specified
	if idx == len(t.ExpectedReflections)-1 && !t.ExpectedLastReflectionHasNoSpace {
		commandReflection = fmt.Sprintf("$ %s ", t.ExpectedReflections[idx])
	}

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the next tab key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", commandReflection)
	return nil
}
