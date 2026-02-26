package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBG1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	typeTestCase := test_cases.CommandResponseTestCase{
		Command:        "type jobs",
		ExpectedOutput: "jobs is a shell builtin",
		SuccessMessage: "Expected type for 'jobs' found",
	}

	if err := typeTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsTestCase := test_cases.JobsBuiltinResponseTestCase{
		SuccessMessage: "✓ Received empty response",
		// Expect no output
		ExpectedOutputItems: []test_cases.JobsBuiltinOutputEntry{},
	}

	if err := jobsTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
