package test_cases

import (
	"fmt"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// PromptTestCaseVT verifies a prompt exists, and that there's no extra output after it.
type PromptTestCaseVT struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	// shouldOmitSuccessLog determines whether a success log should be written.
	//
	// When re-using this test case within other higher-order test cases, emitting success logs
	// all the time can get pretty noisy.
	shouldOmitSuccessLog bool
}

func NewPromptTestCaseVT(expectedPrompt string) PromptTestCaseVT {
	return PromptTestCaseVT{expectedPrompt: expectedPrompt}
}

func NewSilentPromptTestCaseVT(expectedPrompt string) PromptTestCaseVT {
	return PromptTestCaseVT{expectedPrompt: expectedPrompt, shouldOmitSuccessLog: true}
}

func (t PromptTestCaseVT) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	_, err := shell.ReadBytesUntilTimeout(10 * time.Millisecond)
	fullOutput := shell.GetScreenStateForLogging(false)
	actualValue := shell.GetScreenStateSingleRowForLogging(0, false)

	// If reading the prompt failed, or the prompt doesn't match our expectations
	if err != nil || strings.TrimRight(actualValue, "\n") != t.expectedPrompt {
		shell.LogOutput([]byte(fullOutput))

		return fmt.Errorf("Expected prompt (%q) to be printed, got %q", t.expectedPrompt, string(actualValue))
	}

	_, extraOutputErr := shell.ReadBytesUntilTimeout(10 * time.Millisecond)

	extraOutput := shell.GetRowsTillEndForLogging(1, false)

	// Whether the value matches our expectations or not, we print it
	shell.LogOutput([]byte(extraOutput))

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
