package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	testCase1 := test_cases.CDAndPWDTestCase{Directory: "/usr/", Response: "/usr"}
	testCase1.Run(shell, logger)

	testCase2 := test_cases.CDAndPWDTestCase{Directory: "./local/bin", Response: "/usr/local/bin"}
	testCase2.Run(shell, logger)

	testCase3 := test_cases.CDAndPWDTestCase{Directory: "../../", Response: "/usr"}
	testCase3.Run(shell, logger)

	return promptCheck(shell, logger)
}
