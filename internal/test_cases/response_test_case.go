package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ResponseTestCase reads the output from the shell, and verifies that it matches the expected output.
type ResponseTestCase struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	assertion assertions.SingleLineScreenStateAssertion

	shouldOmitSuccessLog bool
}

func NewResponseTestCase(expectedPrompt string, assertion assertions.SingleLineScreenStateAssertion, shouldOmitSuccessLog bool) ResponseTestCase {
	return ResponseTestCase{expectedPrompt: expectedPrompt, assertion: assertion, shouldOmitSuccessLog: shouldOmitSuccessLog}
}

func NewSilentResponseTestCase(expectedPrompt string, assertion assertions.SingleLineScreenStateAssertion) ResponseTestCase {
	return ResponseTestCase{expectedPrompt: expectedPrompt, assertion: assertion, shouldOmitSuccessLog: true}
}

func (t ResponseTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// We can't use the assertion directly as a condition to break out of reads, because of a mismtach in the types that the Read passes and what the assertion expects.
	matchesPromptCondition := func(buf []byte) bool {
		return string(shell_executable.StripANSI(buf)) == t.expectedPrompt
	}

	actualValue, err := shell.ReadBytesUntil(t.assertion.Run())

	if err != nil {
		// If the user sent any output, let's print it before the error message.
		if len(actualValue) > 0 {
			shell.LogOutput(shell_executable.StripANSI(actualValue))
		}

		return fmt.Errorf("Expected prompt (%q) to be printed, got %q", t.expectedPrompt, string(actualValue))
	}

	extraOutput, extraOutputErr := shell.ReadBytesUntilTimeout(10 * time.Millisecond)
	fullOutput := append(actualValue, extraOutput...)

	// Whether the value matches our expectations or not, we print it
	shell.LogOutput(shell_executable.StripANSI(fullOutput))

	// We failed to read extra output
	if extraOutputErr != nil {
		return fmt.Errorf("Error reading output: %v", extraOutputErr)
	}

	if len(extraOutput) > 0 {
		return fmt.Errorf("Found extra output after prompt: %q. (expected just %q)", string(extraOutput), t.expectedPrompt)
	}

	if !t.shouldOmitSuccessLog {
		logger.Successf("✓ Received prompt")
	}

	return nil
}
