package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testP2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
		{CommandType: "wc", CommandName: CUSTOM_WC_COMMAND, CommandMetadata: ""},
		{CommandType: "yes", CommandName: CUSTOM_YES_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	inputs := []string{`yes | head -n 1`, `yes | head -n 10000 | wc`}
	expectedOutputs := []string{"y", fmt.Sprintf("%7d%8d%8d", 10000, 10000, 20000)}

	for i, input := range inputs {
		expectedOutput := expectedOutputs[i]
		testCase := test_cases.CommandResponseTestCase{
			Command:        input,
			ExpectedOutput: expectedOutput,
			SuccessMessage: "✓ Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	return logAndQuit(asserter, nil)
}
