package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/screen_asserter"
)

// CommandResponseTestCase
// Sends a command to the shell
// Verifies that command is printed to the screen `$ <COMMAND>` (we expect the prompt to also be present)
// Reads the output from the shell, and verifies that it matches the expected output
// If any error occurs returns the error from the corresponding assertion
type CommandResponseTestCase struct {
	command string

	// expectedOutput is the expected output string to match against
	expectedOutput string

	// fallbackPatterns is a list of regex patterns to match against
	fallbackPatterns []*regexp.Regexp

	// expectedPatternExplanation is the explanation of the expected pattern to
	// show in the error message in case of failure
	expectedPatternExplanation string
}

func NewCommandResponseTestCase(command string, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) CommandResponseTestCase {
	return CommandResponseTestCase{command: command, expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation}
}

func (t CommandResponseTestCase) Run(screenAsserter *screen_asserter.ScreenAsserter) error {
	err := screenAsserter.Shell.SendCommand(t.command)
	if err != nil {
		return fmt.Errorf("Error sending command: %v", err)
	}

	expectedCommandLine := fmt.Sprintf("$ %s", t.command)
	screenAsserter.PushAssertion(assertions.NewSingleLineScreenStateAssertion(expectedCommandLine, nil, ""))
	screenAsserter.PushAssertion(assertions.NewSingleLineScreenStateAssertion(t.expectedOutput, t.fallbackPatterns, t.expectedPatternExplanation))

	if err := screenAsserter.Shell.ReadUntil(AsBool(screenAsserter.RunWithPromptAssertion)); err != nil {
		if err := screenAsserter.RunWithPromptAssertion(); err != nil {
			return err
		}
	}

	return nil
}

func AsBool(T func() error) func() bool {
	// Takes in a function that takes no params & returns an error
	// Returns the function wrapped in a helper such that it returns a bool
	// in liue of the error, true if the function execution is a success
	return func() bool { return T() == nil }
}
