package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testInvalidCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	invalidCommand := getRandomInvalidCommand()

	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:                    invalidCommand,
		FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(`^(bash: )?` + invalidCommand + `: (command )?not found$`)},
		ExpectedPatternExplanation: invalidCommand + ": command not found",
		SuccessMessage:             "Received command not found message",
	}

	commandReflection := fmt.Sprintf("$ %s", invalidCommand)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("%s: command not found", invalidCommand),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: %s: command not found$`, invalidCommand)),
			regexp.MustCompile(fmt.Sprintf(`^%s: command not found$`, invalidCommand)),
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("✓ Received command not found message")

	return logAndQuit(asserter, nil)
}
