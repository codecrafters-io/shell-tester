package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/creack/pty"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	os.Setenv("PS1", "$ ")
	os.Setenv("BASH_SILENCE_DEPRECATION_WARNING", "1")

	cmd := exec.Command("bash", "--norc", "-i")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	logger := stageHarness.Logger
	command := "nonexistent"
	expectedErrorMessage := fmt.Sprintf("%s: command not found", command)

	a := assertions.BufferAssertion{ExpectedValue: expectedErrorMessage}
	stdBuffer := shell_executable.NewFileBuffer(ptmx)
	stdBuffer.FeedStdin([]byte(command))

	if err := a.Run(&stdBuffer, assertions.CoreTestInexact); err != nil {
		return err
	}
	logger.Debugf("Received message: %q", a.ActualValue)

	if strings.Contains(a.ActualValue, "\n") {
		lines := strings.Split(a.ActualValue, "\n")
		if len(lines) >= 2 {
			a.ActualValue = lines[len(lines)-2]
		}
	}

	logger.Successf("Received error message: %q", a.ActualValue)
	return nil
}
