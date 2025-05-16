package internal

import (
	"os"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH3(stageHarness *test_case_harness.TestCaseHarness) error {
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

	testCase1 := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo hello", ExpectedOutput: "hello"},
			{Command: "echo world", ExpectedOutput: "world"},
			{Command: "pwd", ExpectedOutput: wd},
		},
		LastNCommands: 2,
	}
	if err := testCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	testCase2 := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo foo", ExpectedOutput: "foo"},
			{Command: "echo bar", ExpectedOutput: "bar"},
			{Command: "echo baz", ExpectedOutput: "baz"},
			{Command: "echo CodeCrafters", ExpectedOutput: "CodeCrafters"},
			{Command: "echo is", ExpectedOutput: "is"},
			{Command: "echo great", ExpectedOutput: "great"},
			{Command: "echo foo", ExpectedOutput: "foo"},
			{Command: "pwd", ExpectedOutput: wd},
		},
		LastNCommands: 5,
	}
	if err := testCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
