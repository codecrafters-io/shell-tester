package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	result := NewDetailedResult(executable.ExecutableResult{})
	tries := random.RandomInt(3, 5)
	command := "nonexistent"
	for i := 0; i < tries; i++ {
		b.FeedStdin([]byte(command))

		res, err := b.Result()
		if err != nil {
			return err
		}
		result.UpdateData(res)
		errorMessage := string(result.CurrentCommandStdErr(true))

		if !strings.Contains(errorMessage, command+": command not found") {
			return fmt.Errorf("Expected error message to contain '%s: command not found', but got '%s'", command, errorMessage)
		}
		logger.Successf(strings.Split(errorMessage, "\n")[1])
		if b.HasExited() {
			return fmt.Errorf("Program exited before all commands were sent")
		}
	}

	return nil
}
