package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return logAndQuit(asserter, err)
	}

	invalidCommand := "nonexistent"

	test_case := test_cases.CommandResponseTestCase{
		Command:        invalidCommand,
		ExpectedOutput: fmt.Sprintf("%s: command not found", invalidCommand),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: %s: command not found$`, invalidCommand)),
			regexp.MustCompile(fmt.Sprintf(`^%s: command not found$`, invalidCommand)),
		},
		SuccessMessage: "✓ Received command not found message",
	}

	if err := test_case.Run(asserter, shell, logger); err != nil {
		return logAndQuit(asserter, err)
	}

	return logAndQuit(asserter, nil)
}
