package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPrompt(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := getRandomDirectory()
	if err != nil {
		return err
	}

	// Let's set HOME to a random dir in stage 1 so that CI catches starter code
	// that relies on $HOME being set to a specific dir.
	//
	// We rely on mutating $HOME in the stages that test `cd ~`. Placing this here
	// ensures that we never accept starter code that wouldn't work in those stages.
	shell.Setenv("HOME", randomDir)

	if err := shell.Start(); err != nil {
		return err
	}

	// First prompt assertion
	if err := asserter.Assert(); err != nil {
		return err
	}

	asserter.LogRemainingOutput()
	logger.Successf("✓ Received prompt")

	return nil
}
