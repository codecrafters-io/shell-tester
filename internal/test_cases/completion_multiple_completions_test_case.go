package test_cases

import (
	"fmt"
	"regexp"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/google/shlex"
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

	// ExpectedCompletionOptionsLine is the custom reflection to use
	ExpectedCompletionOptionsLine string

	// If ExpectedCompletion does not match the given completion options
	// the obtained completions are checked against the fallback pattern
	ExpectedCompletionOptionsLineFallbackPatterns []*regexp.Regexp

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

	initialInputReflection := fmt.Sprintf("$ %s", t.RawInput)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: initialInputReflection,
		StayOnSameLine: false,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", initialInputReflection)

	// Send TAB
	for i := range t.TabCount {
		shouldRingBell := i == 0 && t.CheckForBell
		logTabForCompletionOptions(logger, t.ExpectedCompletionOptionsLine, shouldRingBell)

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

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   t.ExpectedCompletionOptionsLine,
		FallbackPatterns: t.ExpectedCompletionOptionsLineFallbackPatterns,
	})

	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Although shlex.Split() is an overkill, it is used, if in the future, we decide to combine
	// quoting stages with completion stages
	words, err := shlex.Split(t.ExpectedCompletionOptionsLine)
	if err != nil {
		return fmt.Errorf("Failed to split expected completion options: %v", err)
	}

	logger.Successf("%d matching completion options found", len(words))

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: initialInputReflection,
		StayOnSameLine: true,
	})

	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Remove autocompletion assertion
	asserter.PopAssertion()

	var assertFuncToRun func() error
	if t.SkipPromptAssertion {
		assertFuncToRun = asserter.AssertWithoutPrompt
	} else {
		assertFuncToRun = asserter.AssertWithPrompt
	}

	if err := assertFuncToRun(); err != nil {
		return err
	}

	// Only log success message when it is provided
	if t.SuccessMessage != "" {
		logger.Successf("%s", t.SuccessMessage)
	}

	return nil
}

func logTabForCompletionOptions(logger *logger.Logger, expectedCompletionOptionsLine string, expectBell bool) {
	if expectBell {
		logger.Infof("Pressed %q (expecting bell to ring)", "<TAB>")
		return
	}

	logger.Infof("Pressed %q (expecting %q as completion options)", "<TAB>", expectedCompletionOptionsLine)
}
