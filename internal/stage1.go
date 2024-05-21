package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPrompt(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	testCase := test_cases.NewSilentPromptTestCase("$ ")

	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	logger.Successf("✓ Received prompt")

	return nil
}
