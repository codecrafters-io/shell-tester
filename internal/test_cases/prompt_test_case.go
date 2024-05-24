package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// PromptTestCase verifies a prompt exists, and that there's no extra output after it.
type PromptTestCase struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	// shouldOmitSuccessLog determines whether a success log should be written.
	//
	// When re-using this test case within other higher-order test cases, emitting success logs
	// all the time can get pretty noisy.
	shouldOmitSuccessLog bool
}

func NewPromptTestCase(expectedPrompt string) PromptTestCase {
	return PromptTestCase{expectedPrompt: expectedPrompt}
}

func NewSilentPromptTestCase(expectedPrompt string) PromptTestCase {
	return PromptTestCase{expectedPrompt: expectedPrompt, shouldOmitSuccessLog: true}
}

func (t PromptTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	matchesPromptCondition := func(buf []byte) bool {
		return string(shell_executable.StripANSI(buf)) == t.expectedPrompt
	}

	actualValue, err := shell.ReadBytesUntil(matchesPromptCondition)

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
