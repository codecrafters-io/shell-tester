package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	command := "nonexistent"
	b.FeedStdin([]byte(command))

	res, err := b.Result()
	if err != nil {
		return err
	}
	result := NewDetailedResult(res)

	errorMessage := string(result.CurrentCommandStdErr(true))

	if !strings.Contains(errorMessage, command+": command not found") {
		return fmt.Errorf("Expected error message to contain '%s: command not found', but got '%s'", command, errorMessage)
	}
	logger.Successf(strings.Split(errorMessage, "\n")[1])

	if b.HasExited() {
		return fmt.Errorf("Program exited before all commands were sent")
	}

	b.FeedStdin([]byte("exit"))
	res, err = b.ResultWithWait()
	if err != nil {
		return err
	}
	result.UpdateData(res)

	if !b.HasExited() {
		return fmt.Errorf("Program did not exit after sending 'exit'")
	}
	logger.Successf(strconv.FormatBool(b.HasExited()))
	logger.Debugf(strconv.Itoa(result.exitCode))

	return nil
}
