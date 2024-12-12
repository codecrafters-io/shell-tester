package internal

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	numberOfCommands := 2

	if err := shell.Start(); err != nil {
		return err
	}

	if err := asserter.Assert(); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		command := "invalid_command_" + strconv.Itoa(i+1)

		shell.SendCommand(command)

		// Command Reflection
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("$ %s", command),
		})

		// TODO: Ensure fallback patterns are accurate, and expected Output is accurate
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("%s: command not found", command),
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`^(bash: )?%s: (command )?not found$`, command)),
			},
		})

		if err := asserter.Assert(); err != nil {
			return err
		}

		// testCase := test_cases.SingleLinePatternMatchTestCase{
		// 	Command:                    command,
		// 	ExpectedPattern:            fmt.Sprintf(`^(bash: )?%s: (command )?not found$`, command),
		// 	ExpectedPatternExplanation: fmt.Sprintf("%s: command not found", command),
		// 	SuccessMessage:             "Received command not found message",
		// }

		// if err := testCase.Run(shell, logger); err != nil {
		// 	return err
		// }
	}

	// TODO: Figure out remaining output in SUCCESS scenario
	asserter.LogRemainingOutput()

	// There must be a prompt after the last command too
	return assertShellIsRunning(shell, logger)
}
