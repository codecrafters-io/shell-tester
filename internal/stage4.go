package internal

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	// We test an inexistent command first, just to make sure the logic works in a "loop"
	testCase := test_cases.RegexTestCase{
		Command:                    "invalid_command_1",
		ExpectedPattern:            regexp.MustCompile(`invalid_command_1: (command )?not found\r\n`),
		ExpectedPatternExplanation: fmt.Sprintf("contain %q", "invalid_command_1: command not found"),
		SuccessMessage:             "Received command not found message",
	}

	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	// We can't use RegexTestCase for the exit command (no output to match on), so we use lower-level methods instead
	promptTestCase := test_cases.NewSilentPromptTestCase("$ ")

	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}

	if err := shell.SendCommand("exit 0"); err != nil {
		return err
	}

	output, readErr := shell.ReadBytesUntilTimeout(1000 * time.Millisecond)
	sanitizedOutput := shell_executable.StripANSI(output)

	// If anything was printed, log it out before we emit error / success logs
	if len(sanitizedOutput) > 0 {
		shell.LogOutput(sanitizedOutput)
	}

	// We're expecting EOF since the program should've terminated
	// HACK: We also allow syscall.EIO since that's what we get on Linux when the process is killed
	// read /dev/ptmx: input/output error *fs.PathError
	if !errors.Is(readErr, io.EOF) && !errors.Is(readErr, syscall.EIO) {
		if readErr == nil {
			return fmt.Errorf("Expected program to exit with 0 exit code, program is still running.")
		} else {
			// TODO: Other than EOF, what other errors could we get? Are they user errors or internal errors?
			return fmt.Errorf("Error reading output: %v", readErr)
		}
	}

	isTerminated, exitCode := shell.WaitForTermination()
	if !isTerminated {
		return fmt.Errorf("Expected program to exit with 0 exit code, program is still running.")
	}

	logger.Successf("✓ Program exited successfully")

	if exitCode != 0 {
		return fmt.Errorf("Expected 0 as exit code, got %d", exitCode)
	}

	// Most shells return nothing but bash returns the string "exit" when it exits, we allow both styles
	if len(sanitizedOutput) > 0 && strings.TrimSpace(string(sanitizedOutput)) != "exit" {
		return fmt.Errorf("Expected no output after exit command, got %q", string(sanitizedOutput))
	}

	logger.Successf("✓ No output after exit command")

	return nil
}
