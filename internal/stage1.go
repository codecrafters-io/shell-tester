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

	// XXX: selector ?
	// XXX: Why is stdout empty ?
	buffer, err := b.ReadBuffer("stderr")
	if err != nil {
		return err
	}

	cleanedBuffer := removeControlSequence(buffer)
	prompt := string(cleanedBuffer)
	expectedPrompt := "$ "

	if len(prompt) == 0 {
		return fmt.Errorf("Expected to receive prompt, but got nothing")
	}

	if !strings.EqualFold(prompt, expectedPrompt) {
		return fmt.Errorf("Expected prompt to be %q, but got %q", expectedPrompt, prompt)
	}
	logger.Successf("Received prompt: %q", prompt)

	if b.HasExited() {
		return fmt.Errorf("Expected shell to be running, but it has exited")
	}
	logger.Successf("Shell is still running")

	return nil
}
