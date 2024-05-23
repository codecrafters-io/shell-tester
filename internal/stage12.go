package internal

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	testCase1 := test_cases.CDAndPWDTestCase{Directory: "/usr/local/bin", Response: "/usr/local/bin"}
	testCase1.Run(shell, logger)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Error getting home directory: %v", err)
	}

	testCase2 := test_cases.CDAndPWDTestCase{Directory: "~", Response: homeDir}
	testCase2.Run(shell, logger)

	return assertShellIsRunning(shell, logger)
}