package internal

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const commandSuccessMessage = "✓ Command executed successfully"

func testH2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	randomWords1 := strings.Join(random.RandomWords(2), " ")
	randomWords2 := strings.Join(random.RandomWords(2), " ")
	randomWords3 := strings.Join(random.RandomWords(2), " ")

	previousCommandsTestCases := []test_cases.CommandResponseTestCase{
		{Command: "echo " + randomWords1, ExpectedOutput: randomWords1, SuccessMessage: commandSuccessMessage},
		{Command: "echo " + randomWords2, ExpectedOutput: randomWords2, SuccessMessage: commandSuccessMessage},
		{Command: "echo " + randomWords3, ExpectedOutput: randomWords3, SuccessMessage: commandSuccessMessage},
	}
	for _, command := range previousCommandsTestCases {
		if err := command.Run(asserter, shell, logger); err != nil {
			return err
		}
	}
	previousCommands := []string{}
	for _, command := range previousCommandsTestCases {
		previousCommands = append(previousCommands, command.Command)
	}

	testCase := test_cases.HistoryTestCase{
		PreviousCommands: previousCommands,
		SuccessMessage:   "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
