package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBG3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Launch a program in the background
	backgroundLaunchCommand := "sleep 100"
	backgroundLaunchTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           backgroundLaunchCommand,
		SuccessMessage:    "✓ Output includes job number with PID",
		ExpectedJobNumber: 1,
	}

	if err := backgroundLaunchTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Assert the job output
	jobsTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{{
			JobNumber:     1,
			Status:        "Running",
			LaunchCommand: backgroundLaunchCommand,
			Marker:        test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ 1 entry matches the running job",
	}

	if err := jobsTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
