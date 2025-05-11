package internal

import (
	"os"

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

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	testCase := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo hello", ExpectedOutput: "hello"},
			{Command: "echo world", ExpectedOutput: "world"},
			{Command: "pwd", ExpectedOutput: wd},
		},
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
