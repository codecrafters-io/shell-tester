package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	typeTestCase := test_cases.TypeOfCommandTestCase{
		Command: "complete",
	}

	if err := typeTestCase.RunForBuiltin(asserter, shell, logger); err != nil {
		return err
	}

	completeTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        "complete",
		SuccessMessage: "✓ No output found",
	}

	if err := completeTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
