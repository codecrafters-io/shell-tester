package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	testCase := test_cases.InvalidCommandTestCase{
		Command: getRandomInvalidCommand(),
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	refTestCase := test_cases.ExitTestCase{
		Command:          "exit 0",
		ExpectedExitCode: 0,
	}
	if err := refTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	logger.Successf("✓ No output after exit command")

	return logAndQuit(asserter, nil)
}
