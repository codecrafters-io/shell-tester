package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testP5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "ls", CommandName: CUSTOM_LS_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	input := `ls | type exit`
	expectedOutput := `exit is a shell builtin`

	singleLineTestCase := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
