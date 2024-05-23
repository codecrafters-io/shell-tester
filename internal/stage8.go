package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/elf_executable"
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

	randomString := elf_executable.GetRandomString()
	err := elf_executable.CreateELFexecutable(randomString+"\n", "my_exe")
	if err != nil {
		return err
	}

	commandsWithResponse := []string{
		"./my_exe",
	}

	for _, command := range commandsWithResponse {
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(randomString + "\r\n"),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", randomString),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
