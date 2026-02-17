package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type InputAndCompletionPair struct {
	Input              string
	ExpectedCompletion string
}

// PartialCompletionsTestCase is a test case that does the following:
// For each InputAndCompletionPair it performs the following steps:
// It sends the input to the shell
// Asserts the typed input to appear on the TTY
// Sends TAB and expects it to autocomplete to the completion text
// Repeats this process for all the completion pairs in the same prompt line
type PartialCompletionsTestCase struct {
	// InputAndCompletionPairs should specify the text to type and the expected completion
	// In each step of this test case
	InputAndCompletionPairs []InputAndCompletionPair

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t PartialCompletionsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Send the given input from the input and completion pair
	// Test for the input appearance
	// Send tab
	// Expect the completion to be made
	for idx := 0; idx < len(t.InputAndCompletionPairs); idx++ {
		// Starting from the second input and expected completion pair,
		// the expected completion from the previous pair will have already been present in the line already
		preExistingInputOnLine := ""
		if idx > 0 {
			preExistingInputOnLine = t.InputAndCompletionPairs[idx-1].ExpectedCompletion
		}

		if err := t.runInputAppearanceAssertion(asserter, shell, logger, t.InputAndCompletionPairs[idx], preExistingInputOnLine); err != nil {
			return err
		}

		if err := t.runTabCompletionAssertion(asserter, shell, logger, t.InputAndCompletionPairs[idx]); err != nil {
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

func (t PartialCompletionsTestCase) runInputAppearanceAssertion(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, inputCompletionPair InputAndCompletionPair, preExistingInputOnLine string) error {
	// We don't need to run the test case if the input is empty
	if inputCompletionPair.Input == "" {
		return nil
	}

	// Log the details of the text before sending it
	logTypedText(logger, inputCompletionPair.Input)

	// Send the text to the shell
	if err := shell.SendText(inputCompletionPair.Input); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	expectedPromptLine := fmt.Sprintf("$ %s", preExistingInputOnLine+inputCompletionPair.Input)

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: expectedPromptLine,
		StayOnSameLine: true,
	})

	// Run the assertion, before sending the tab key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	asserter.PopAssertion()
	logger.Successf("✓ Prompt line matches %q", expectedPromptLine)
	return nil
}

func (t PartialCompletionsTestCase) runTabCompletionAssertion(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, inputCompletionPair InputAndCompletionPair) error {
	// Send TAB
	logTab(logger, inputCompletionPair.ExpectedCompletion, false)
	if err := shell.SendText("\t"); err != nil {
		return fmt.Errorf("Error sending text to shell: %v", err)
	}

	expectedCompletion := fmt.Sprintf("$ %s", inputCompletionPair.ExpectedCompletion)

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: expectedCompletion,
		StayOnSameLine: true,
	})

	// Run the assertion, before sending the next tab key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	asserter.PopAssertion()

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", expectedCompletion)
	return nil
}
