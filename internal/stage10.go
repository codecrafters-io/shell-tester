package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	directory, err := GetRandomDirectory()
	if err != nil {
		return err
	}

	testCase := test_cases.CDAndPWDTestCase{Directory: directory, Response: directory}
	err = testCase.Run(shell, logger)
	if err != nil {
		return err
	}

	directory = "/non-existing-directory"
	command := fmt.Sprintf("cd %s", directory)

	failureTestCase := test_cases.RegexTestCase{
		Command:                    command,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`^(can't cd to %s|(bash: cd: )?%s: No such file or directory)\r\n`, directory, directory)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", fmt.Sprintf(`cd: %s: No such file or directory\n`, directory)),
		SuccessMessage:             "Received error message",
	}

	if err := failureTestCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
