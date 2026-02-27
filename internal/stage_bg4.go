package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/dustin/go-humanize/english"
)

func testBG4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	commands := []string{"sleep 100", "sleep 200", "sleep 300"}

	if err := launchBgCommandAndAssertJobs(asserter, shell, logger, commands); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

// launchBgCommandAndAssertJobs launches the given bgCommands one by one
// with a 'jobs' call after each launch
func launchBgCommandAndAssertJobs(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, bgCommands []string) error {
	type jobInfo struct {
		JobNumber int
		Command   string
	}

	var jobs []jobInfo

	for i, bgCommand := range bgCommands {
		backgroundLaunchTestCase := test_cases.BackgroundCommandResponseTestCase{
			Command:           bgCommand,
			SuccessMessage:    "✓ Output includes job number with PID",
			ExpectedJobNumber: i + 1,
		}

		if err := backgroundLaunchTestCase.Run(asserter, shell, logger); err != nil {
			return err
		}

		jobs = append(jobs, jobInfo{JobNumber: i + 1, Command: bgCommand})

		jobsOutputEntries := make([]test_cases.BackgroundJobStatusEntry, 0, len(jobs))

		for i, job := range jobs {
			// Default marker is unmarked
			marker := test_cases.UnmarkedJob

			// If the job was recently launched, it is the 'Current' job
			if i == len(jobs)-1 {
				marker = test_cases.CurrentJob
				// If the job was launched previously, it is the 'Previous' job
			} else if i == len(jobs)-2 {
				marker = test_cases.PreviousJob
			}

			jobsOutputEntries = append(jobsOutputEntries, test_cases.BackgroundJobStatusEntry{
				JobNumber:     job.JobNumber,
				Status:        "Running",
				LaunchCommand: job.Command,
				Marker:        marker,
			})
		}

		jobsTestCase := test_cases.JobsBuiltinResponseTestCase{
			ExpectedOutputEntries: jobsOutputEntries,
			SuccessMessage: fmt.Sprintf(
				"✓ %s match the running %s",
				english.Plural(len(jobsOutputEntries), "entry", "entries"),
				english.Plural(len(jobsOutputEntries), "job", "jobs"),
			),
		}

		if err := jobsTestCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}
	return nil
}
