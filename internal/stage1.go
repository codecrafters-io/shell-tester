package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
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
	expectedPrompt := "$ "

	a := assertions.ExactBufferAssertion{ExpectedValue: expectedPrompt}
	if err := a.Run(b, "stderr"); err != nil {
		return err
	}
	logger.Successf("Received prompt: %q", a.ActualValue)

	if b.HasExited() {
		return fmt.Errorf("Expected shell to be running, but it has exited")
	}
	logger.Successf("Shell is still running")

	return nil
}
