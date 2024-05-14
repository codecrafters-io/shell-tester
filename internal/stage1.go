package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPrompt(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	b.FeedStdin([]byte(""))

	res, err := b.Result()
	if err != nil {
		return err
	}
	result := NewDetailedResult(res)

	// XXX: Why is stdout empty ?
	prompt := strings.TrimSpace(string(result.CurrentCommandStdErr(false)))

	if len(prompt) == 0 {
		return fmt.Errorf("Expected to receive prompt, but got nothing")
	}

	// bash will send extra characters apart from $ prompt
	if !strings.Contains(prompt, "$") {
		return fmt.Errorf("Expected prompt to be '$', but got '%s'", prompt)
	}
	logger.Successf("Received prompt: %q", strings.Split(prompt, "\n")[1])

	if b.HasExited() {
		return fmt.Errorf("Program has exited")
	}
	logger.Successf("Shell is still running")

	return nil
}
