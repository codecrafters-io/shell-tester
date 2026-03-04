package internal

import (
	"fmt"
	"path/filepath"
	"time"

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

	fifoBaseNames := random.RandomWords(2)
	fifoPath1 := filepath.Join(
		"/tmp",
		fmt.Sprintf("%s-%d", fifoBaseNames[0], random.RandomInt(1, 100)),
	)
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath1, 0644); err != nil {
		return err
	}

	fifoPath2 := filepath.Join(
		"/tmp",
		fmt.Sprintf("%s-%d", fifoBaseNames[1], random.RandomInt(1, 100)),
	)

	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath2, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Launch "sleep 500"
	sleepCommand := "sleep 500"
	bgSleepTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepCommand,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Output includes job number with PID",
	}
	if err := bgSleepTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Launch 'cat read pattern' to hang the process indefinitely
	bgCatCommand1 := fmt.Sprintf("cat %s", fifoPath1)
	bgCatTestCase1 := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgCatCommand1,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Output includes job number with PID",
	}
	if err := bgCatTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Launch cat that reads from the fifo
	bgCatCommand2 := fmt.Sprintf("cat %s", fifoPath2)
	bgCatTestCase2 := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgCatCommand2,
		ExpectedJobNumber: 3,
		SuccessMessage:    "✓ Output includes job number with PID",
	}

	if err := bgCatTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write an empty string to the first fifo
	if err := WriteToFile(stageHarness, fifoPath1, ""); err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	// Call jobs for the first time
	jobsBuiltinTestCase1 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.UnmarkedJob},
			{JobNumber: 2, Status: "Done", LaunchCommand: bgCatCommand1, Marker: test_cases.PreviousJob},
			{JobNumber: 3, Status: "Running", LaunchCommand: bgCatCommand2, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 3 entries in the output",
	}

	if err := jobsBuiltinTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write empty string to the second fifo
	if err := WriteToFile(stageHarness, fifoPath2, ""); err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	// Call jobs for the second time
	jobsBuiltinTestCase2 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.PreviousJob},
			{JobNumber: 3, Status: "Done", LaunchCommand: bgCatCommand2, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 2 entries in the output",
	}
	if err := jobsBuiltinTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Call jobs again
	jobsBuiltinTestCase3 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ 1 entry matches the running job",
	}
	if err := jobsBuiltinTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
