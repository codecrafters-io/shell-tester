package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

	invalidCommand := getRandomInvalidCommand()

	// We test an invalid command first, just to make sure the logic works in a "loop"
	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:                    invalidCommand,
		FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(`^(bash: )?` + invalidCommand + `: (command )?not found$`)},
		ExpectedPatternExplanation: invalidCommand + ": command not found",
		SuccessMessage:             "Received command not found message",
	}

	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	// We can't use SingleLineOutputTestCase for the exit command (no output to match on), so we use lower-level methods instead
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
	if !errors.Is(readErr, shell_executable.ErrProgramExited) {
		if readErr == nil {
			return fmt.Errorf("Expected program to exit with 0 exit code, program is still running.")
		} else {
			// TODO: Other than ErrProgramExited, what other errors could we get? Are they user errors or internal errors?
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
