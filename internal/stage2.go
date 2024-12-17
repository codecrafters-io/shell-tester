package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
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

	// We are seperating this out because we don't want to assert 
	// The prompt at the end
	testCase := test_cases.CommandReflectionTestCase{
		Command: invalidCommand,
	}
	if err := testCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("%s: command not found", invalidCommand),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: %s: command not found$`, invalidCommand)),
			regexp.MustCompile(fmt.Sprintf(`^%s: command not found$`, invalidCommand)),
		},
	})

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	logger.Successf("✓ Received command not found message")

	return logAndQuit(asserter, nil)
}
