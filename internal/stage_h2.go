package internal

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	randomWords1 := strings.Join(random.RandomWords(2), " ")
	randomWords2 := strings.Join(random.RandomWords(2), " ")
	randomWords3 := strings.Join(random.RandomWords(2), " ")

	testCase := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo " + randomWords1, ExpectedOutput: randomWords1},
			{Command: "echo " + randomWords2, ExpectedOutput: randomWords2},
			{Command: "echo " + randomWords3, ExpectedOutput: randomWords3},
		},
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
