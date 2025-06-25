package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	tmpHomeDir, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}
	shell.Setenv("HOME", tmpHomeDir)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	directory, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	cdTestCase1 := test_cases.CDAndPWDTestCase{Directory: directory, Response: directory}
	err = cdTestCase1.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	echoTestCase := test_cases.CommandResponseTestCase{
		Command:          "echo ~",
		ExpectedOutput:   tmpHomeDir,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected response",
	}
	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	cdTestCase2 := test_cases.CDAndPWDTestCase{Directory: "~", Response: tmpHomeDir}
	err = cdTestCase2.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
