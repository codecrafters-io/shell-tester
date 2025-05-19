package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	testCase := test_cases.HistoryTestCase{
		SuccessMessage: "✓ History command works as expected",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "ls dist/", ExpectedOutput: "main.out"},
			{Command: "cd dist/", ExpectedOutput: ""},
			{Command: "ls", ExpectedOutput: "main.out"},
		},
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
