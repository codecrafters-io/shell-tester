package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	builtIns := []string{"echo", "exit", "type"}

	if err := shell.Start(); err != nil {
		return err
	}

	for _, builtIn := range builtIns {
		command := fmt.Sprintf("type %s", builtIn)

		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`%s is a( special)? shell builtin\r\n`, builtIn)),
			ExpectedPatternExplanation: fmt.Sprintf("match %q\n", fmt.Sprintf(`%s is a shell builtin`, builtIn)),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	invalidCommands := []string{"nonexistent", "nonexistentcommand"}

	for _, invalidCommand := range invalidCommands {
		command := fmt.Sprintf("type %s", invalidCommand)

		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`(bash: type: )?%s[:]? not found\r\n`, invalidCommand)),
			ExpectedPatternExplanation: fmt.Sprintf("contain %q", fmt.Sprintf(`%s not found\n`, invalidCommand)),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
