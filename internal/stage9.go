package internal

import (
	"fmt"
	"os"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testpwd(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	command := "pwd"
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Error getting cwd: %v", err)
	}

	testCase := test_cases.RegexTestCase{
		Command:                    command,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`%s\r\n`, cwd)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", cwd),
		SuccessMessage:             "Received current working directory response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	// ToDo: Add check for shell still running

	return nil
}
