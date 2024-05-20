package internal

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	if err := shell.AssertPrompt("$ "); err != nil {
		return err
	}

	if err := shell.SendCommand("inexistent"); err != nil {
		return err
	}

	if err := shell.AssertOutputMatchesRegex(regexp.MustCompile(`inexistent: (command )?not found\r\n`)); err != nil {
		return err
	}

	logger.Successf("✓ Received command not found message")

	// At this stage the user might or might not have implemented a REPL to print the prompt again, so we won't test further

	return nil
}
