package internal

import (
	"os"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	tmpHomeDir, err := getRandomDirectory()
	if err != nil {
		return err
	}
	os.Setenv("HOME", tmpHomeDir)

	if err := shell.Start(); err != nil {
		return err
	}

	directory, err := getRandomDirectory()
	if err != nil {
		return err
	}

	testCase1 := test_cases.CDAndPWDTestCase{Directory: directory, Response: directory}
	err = testCase1.Run(shell, logger)
	if err != nil {
		return err
	}

	testCase2 := test_cases.CDAndPWDTestCase{Directory: "~", Response: tmpHomeDir}
	err = testCase2.Run(shell, logger)
	if err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
