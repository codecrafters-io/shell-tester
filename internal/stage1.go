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

	expectedPrompt := "$ "

	a := assertions.ExactBufferAssertion{ExpectedValue: expectedPrompt}
	truncatedStdErrBuf := shell_executable.NewTruncatedBuffer(b.GetStdErrBuffer())

	if err := a.Run(&truncatedStdErrBuf); err != nil {
		return err
	}
	logger.Successf("Received prompt: %q", a.ActualValue)

	if b.HasExited() {
		return fmt.Errorf("Expected shell to be running, but it has exited")
	}
	logger.Successf("Shell is still running")

	return nil
}
