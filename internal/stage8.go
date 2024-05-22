package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRun(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	promptTestCase := test_cases.NewPromptTestCase("$ ")
	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}

	fileContent := "Hello, World!"
	commandsWithNoReply := []string{
		"rm -rf /tmp/foo",
		"mkdir -p /tmp/foo/bar/baz",
		"touch /tmp/foo/bar/baz/file.txt",
		fmt.Sprintf("echo '%s' > /tmp/foo/bar/baz/file.txt", fileContent),
	}

	for i, command := range commandsWithNoReply {
		testCase := test_cases.EmptyResponseTestCase{
			Command:        command,
			SuccessMessage: "Received empty response",
		}
		if err := testCase.Run(shell, logger, i == len(commandsWithNoReply)-1); err != nil {
			return err
		}
	}

	commandsWithResponse := []string{
		"cat /tmp/foo/bar/baz/file.txt",
	}

	for _, command := range commandsWithResponse {
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(fileContent + "\r\n"),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", fileContent),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	// ToDo: Add check for shell still running

	return nil
}
