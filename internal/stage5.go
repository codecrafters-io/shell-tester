package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEcho(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	numberOfCommands := random.RandomInt(2, 4)

	if err := shell.Start(); err != nil {
		return err
	}

	// First prompt assertion
	if err := asserter.Assert(); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		words := strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
		command := fmt.Sprintf("echo %s", words)

		testCase := test_cases.CommandResponseTestCase{
			Command:          command,
			ExpectedOutput:   words,
			FallbackPatterns: nil,
			SuccessMessage:   "Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	asserter.LogRemainingOutput()
	return nil
}
