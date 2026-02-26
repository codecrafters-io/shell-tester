package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBG7(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	fifoPath := fmt.Sprintf("/tmp/%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Spawn background process: sleep 500
	sleepCommand := "sleep 500"
	bgSleepTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepCommand,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Received entry for the launched job",
	}
	if err := bgSleepTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Grep read pattern
	grepPattern := random.RandomWord()
	bgGrepCommand := fmt.Sprintf("grep -q %s %s", grepPattern, fifoPath)
	bgGrepTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Received entry for the launched job",
	}
	if err := bgGrepTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Run jobs
	jobsBuiltinTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.PreviousJob},
			{JobNumber: 2, Status: "Running", LaunchCommand: bgGrepCommand, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 2 entries for the running jobs",
	}
	if err := jobsBuiltinTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to fifo
	if err := WriteToFile(stageHarness, fifoPath, grepPattern); err != nil {
		return err
	}

	// Issue an echo command and expect the reaped job entry will follow the echoed text
	echoArgument := random.RandomWord()
	echoTestCase := test_cases.CommandResponseWithReapedJobsTestCase{
		Command:               fmt.Sprintf("echo %s", echoArgument),
		ExpectedCommandOutput: echoArgument,
		ExpectedReapedJobEntries: []*test_cases.BackgroundJobStatusEntry{{
			JobNumber:     2,
			Status:        "Done",
			LaunchCommand: bgGrepCommand,
			Marker:        test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ Received command output followed by an entry for the reaped job",
	}

	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Call jobs — only sleep (job 1) remains
	jobsBuiltinTestCase2 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{{
			JobNumber:     1,
			Status:        "Running",
			LaunchCommand: sleepCommand,
			Marker:        test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ Received 1 entry for the remaining running job",
	}
	if err := jobsBuiltinTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
