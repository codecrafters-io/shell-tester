package internal

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	randomWords1 := strings.Join(random.RandomWords(2), " ")
	randomWords2 := strings.Join(random.RandomWords(2), " ")
	randomWords3 := strings.Join(random.RandomWords(2), " ")

	testCase1 := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo " + randomWords1, ExpectedOutput: randomWords1},
			{Command: "echo " + randomWords2, ExpectedOutput: randomWords2},
			{Command: "echo " + randomWords3, ExpectedOutput: randomWords3},
		},
		LastNCommands: 2,
	}
	if err := testCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	randomWords4 := strings.Join(random.RandomWords(2), " ")
	randomWords5 := strings.Join(random.RandomWords(2), " ")
	randomWords6 := strings.Join(random.RandomWords(2), " ")
	randomWords7 := strings.Join(random.RandomWords(2), " ")
	randomWords8 := strings.Join(random.RandomWords(2), " ")
	randomWords9 := strings.Join(random.RandomWords(2), " ")
	randomWords10 := strings.Join(random.RandomWords(2), " ")
	randomWords11 := strings.Join(random.RandomWords(2), " ")

	testCase2 := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo " + randomWords4, ExpectedOutput: randomWords4},
			{Command: "echo " + randomWords5, ExpectedOutput: randomWords5},
			{Command: "echo " + randomWords6, ExpectedOutput: randomWords6},
			{Command: "echo " + randomWords7, ExpectedOutput: randomWords7},
			{Command: "echo " + randomWords8, ExpectedOutput: randomWords8},
			{Command: "echo " + randomWords9, ExpectedOutput: randomWords9},
			{Command: "echo " + randomWords10, ExpectedOutput: randomWords10},
			{Command: "echo " + randomWords11, ExpectedOutput: randomWords11},
		},
		LastNCommands: random.RandomInt(3, 6),
	}
	if err := testCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Can remove, just did it a second time to allow more insight with varying number `n`
	testCase3 := test_cases.HistoryTestCase{
		SuccessMessage: "✓ Received expected response",
		CommandsBeforeHistory: []test_cases.CommandOutputPair{
			{Command: "echo " + randomWords4, ExpectedOutput: randomWords4},
			{Command: "echo " + randomWords5, ExpectedOutput: randomWords5},
			{Command: "echo " + randomWords6, ExpectedOutput: randomWords6},
			{Command: "echo " + randomWords7, ExpectedOutput: randomWords7},
			{Command: "echo " + randomWords8, ExpectedOutput: randomWords8},
			{Command: "echo " + randomWords9, ExpectedOutput: randomWords9},
			{Command: "echo " + randomWords10, ExpectedOutput: randomWords10},
			{Command: "echo " + randomWords11, ExpectedOutput: randomWords11},
		},
		LastNCommands: random.RandomInt(5, 7),
	}
	if err := testCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
