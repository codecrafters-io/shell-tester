package internal

import (
	"fmt"
	"os"
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

	directory := "/usr/local/bin"
	command := fmt.Sprintf("cd %s", directory)
	_, err := os.Stat(directory)
	if err != nil {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error creating tmp directory: %v", err)
		}
	}

	promptTestCase := test_cases.NewPromptTestCase("$ ")
	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}

	if err := shell.SendCommand(command); err != nil {
		return err
	}

	command = "pwd"

	testCase := test_cases.RegexTestCase{
		Command:                    command,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`%s\r\n`, directory)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", directory),
		SuccessMessage:             "Received current working directory response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	directory = "/non-existing-directory"
	command = fmt.Sprintf("cd %s", directory)

	testCase = test_cases.RegexTestCase{
		Command:                    command,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`(can't cd to %s|%s: No such file or directory)\r\n`, directory, directory)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", fmt.Sprintf(`%s: No such file or directory\r\n`, directory)),
		SuccessMessage:             "Received error message",
	}

	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return promptCheck(shell, logger)
}
