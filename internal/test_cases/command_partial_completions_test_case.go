package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandPartialCompletionsTestCase is a test case that:
// Sends a command to the shell
// Asserts that the prompt line reflects the command
// for each partial auto-completion:
// Sends TAB
// Asserts that the expected reflection is printed to the screen
// And sends the subsequent input
// If any error occurs returns the error from the corresponding assertion
type CommandPartialCompletionsTestCase struct {
	// Inputs is the list of inputs to send to the shell
	// They are send one by one, interleaved with TABs
	// The shell is expected to auto-complete expected reflections
	Inputs []string

	// ExpectedReflections is the list of expected reflections to use
	ExpectedReflections []string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandPartialCompletionsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if len(t.Inputs) != len(t.ExpectedReflections) {
		panic("Inputs and ExpectedReflections must have the same length")
	}

	// The entire flow is repeated for each input & expected reflection
	for idx := 0; idx < len(t.ExpectedReflections); idx++ {
		// Log the details of the command before sending it
		if t.Inputs[idx] != "" {
			logCommand(logger, t.Inputs[idx])
		}

		// Send the command to the shell
		if err := shell.SendCommandRaw(t.Inputs[idx]); err != nil {
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

		// Send TAB
		logTab(logger, t.ExpectedReflections[idx], false)
		if err := shell.SendCommandRaw("\t"); err != nil {
			return fmt.Errorf("Error sending command to shell: %v", err)
		}

		// For all partial auto-completions, we expect *NO* space at the end
		commandReflection := fmt.Sprintf("$ %s", t.ExpectedReflections[idx])
		// For the last auto-completion, we expect a space at the end
		if idx == len(t.ExpectedReflections)-1 {
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
