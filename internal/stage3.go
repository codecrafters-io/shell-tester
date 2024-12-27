package internal

import (
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	numberOfCommands := random.RandomInt(3, 6)

	for i := 0; i < numberOfCommands; i++ {
		testCase := test_cases.InvalidCommandTestCase{
			Command: "invalid_command_" + strconv.Itoa(i+1),
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// AssertWithPrompt() already makes sure that the prompt is present in the last row
	return logAndQuit(asserter, nil)
}
