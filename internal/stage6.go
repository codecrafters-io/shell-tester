package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	builtIns := []string{"echo", "exit", "type"}

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	for _, builtIn := range builtIns {
		command := fmt.Sprintf("type %s", builtIn)

		testCase := test_cases.CommandResponseTestCase{
			Command:          command,
			ExpectedOutput:   fmt.Sprintf(`%s is a shell builtin`, builtIn),
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^%s is a( special)? shell builtin$`, builtIn))},
			SuccessMessage:   "✓ Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	invalidCommands := getRandomInvalidCommands(2)

	for _, invalidCommand := range invalidCommands {
		command := fmt.Sprintf("type %s", invalidCommand)

		testCase := test_cases.CommandResponseTestCase{
			Command:          command,
			ExpectedOutput:   fmt.Sprintf("%s: not found", invalidCommand),
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(bash: type: )?%s[:]? not found$`, invalidCommand))},
			SuccessMessage:   "✓ Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	return logAndQuit(asserter, nil)
}
