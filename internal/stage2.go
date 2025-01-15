package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testInvalidCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	testCase := test_cases.InvalidCommandTestCase{
		Command: "invalid",
	}
	if err := testCase.RunWithoutNextPromptAssertion(asserter, shell, logger); err != nil {
		return err
	}

	logger.Successf("✓ Received command not found message")

	// TODO: Printing this prompt makes fixture tests flaky
	return nil
}
