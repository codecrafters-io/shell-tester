package internal

import (
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	directory, err := getRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	separator := os.PathSeparator
	parentDirs := strings.Split(directory, string(separator))

	// first 2 dirs, /tmp/foo -> /tmp/foo
	dir := string(separator) + path.Join(parentDirs[:len(parentDirs)-2]...)
	testCase1 := test_cases.CDAndPWDTestCase{Directory: dir, Response: dir}
	err = testCase1.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	// go deeper, ./bar/baz -> /tmp/foo/bar/baz
	dir = "." + string(separator) + path.Join(parentDirs[len(parentDirs)-2:]...)
	absoluteDir := string(separator) + path.Join(parentDirs...)
	testCase2 := test_cases.CDAndPWDTestCase{Directory: dir, Response: absoluteDir}
	err = testCase2.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	// go back, ../../../ -> /tmp
	absoluteDir = string(separator) + path.Join(parentDirs[:len(parentDirs)-3]...)
	testCase3 := test_cases.CDAndPWDTestCase{Directory: "../../../", Response: absoluteDir}
	err = testCase3.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
