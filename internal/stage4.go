package internal

import (
	"fmt"
	"io"
	"regexp"
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
		Command:                    "inexistent",
		ExpectedPattern:            regexp.MustCompile(`inexistent: (command )?not found\r\n`),
		ExpectedPatternExplanation: fmt.Sprintf("contain %q", "inexistent: command not found"),
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

	// If anything was printed, log it out before we emit error / success logs
	if len(output) > 0 {
		shell.LogOutput(shell_executable.StripANSI(output))
	}

	// Either the err is "nil", so we didn't reach EOF
	if readErr == nil {
		return fmt.Errorf("Expected program to terminate, program is still running.")
	} else if readErr != io.EOF {
		// TODO: Other than EOF, what other errors could we get? Are they user errors or internal errors?
		return fmt.Errorf("Error reading output: %v", readErr)
	}

	// If the err IS io.EOF, the program exited successfully
	logger.Successf("✓ Program exited successfully")

	return nil
}
