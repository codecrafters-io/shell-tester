package internal

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	numberOfCommands := random.RandomInt(3, 6)

	if err := shell.Start(); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		command := "invalid_command_" + strconv.Itoa(i+1)

		testCase := test_cases.SingleLineOutputTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`^(bash: )?%s: (command )?not found$`, command)),
			ExpectedPatternExplanation: fmt.Sprintf("contain %q", fmt.Sprintf("%s: command not found\n", command)),
			SuccessMessage:             "Received command not found message",
		}

		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	// There must be a prompt after the last command too
	return assertShellIsRunning(shell, logger)
}
