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

	if err := shell.Start(); err != nil {
		return err
	}

	// First prompt assertion
	if err := asserter.Assert(); err != nil {
		return err
	}

	invalidCommand := "nonexistent"

	commandResponseTestCase := test_cases.CommandResponseTestCase{
		Command:        invalidCommand,
		ExpectedOutput: fmt.Sprintf("%s: command not found", invalidCommand),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: %s: command not found$`, invalidCommand)),
			regexp.MustCompile(fmt.Sprintf(`^%s: command not found$`, invalidCommand)),
		},
		SuccessMessage: "✓ Received command not found message",
	}

	if err := commandResponseTestCase.Run(shell, logger, asserter); err != nil {
		return err
	}

	asserter.LogRemainingOutput()
	return nil
}
