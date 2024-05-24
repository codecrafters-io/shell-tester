package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

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

	directory, err := getRandomDirectory()
	if err != nil {
		return err
	}

	separator := os.PathSeparator
	parentDirs := strings.Split(directory, string(separator))
	fmt.Println(directory, separator, parentDirs, len(parentDirs))

	dir := string(separator) + path.Join(parentDirs[:len(parentDirs)-2]...)
	testCase1 := test_cases.CDAndPWDTestCase{Directory: dir, Response: dir}
	err = testCase1.Run(shell, logger)
	if err != nil {
		return err
	}

	dir = "." + string(separator) + path.Join(parentDirs[len(parentDirs)-2:]...)
	absoluteDir := string(separator) + path.Join(parentDirs...)
	testCase2 := test_cases.CDAndPWDTestCase{Directory: dir, Response: absoluteDir}
	err = testCase2.Run(shell, logger)
	if err != nil {
		return err
	}

	absoluteDir = string(separator) + path.Join(parentDirs[:len(parentDirs)-3]...)
	testCase3 := test_cases.CDAndPWDTestCase{Directory: "../../../", Response: absoluteDir}
	err = testCase3.Run(shell, logger)
	if err != nil {
		return err
	}
	return assertShellIsRunning(shell, logger)
}
