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

func testBG9(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	if err := testBg9ResetToZero(stageHarness); err != nil {
		return err
	}

	logger.Infof("Tearing down shell")

	return testBg9Recycle(stageHarness)
}

func testBg9ResetToZero(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	logger := stageHarness.Logger

	fifoPath1 := filepath.Join(
		"/tmp",
		fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100)),
	)
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath1, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	bgCatCommand1 := fmt.Sprintf("cat %s", fifoPath1)
	bgCatTestCase := &test_cases.BackgroundCommandResponseTestCase{
		Command:           bgCatCommand1,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Output includes job number with PID",
	}

	if err := bgCatTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write empty string to FIFO to make job 1 'Done'
	if err := WriteToFile(stageHarness, fifoPath1, ""); err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	echoArgs := random.RandomWords(2)
	echoTestCase := test_cases.CommandResponseWithReapedJobsTestCase{
		Command:               fmt.Sprintf("echo %s", echoArgs[0]),
		ExpectedCommandOutput: echoArgs[0],
		ExpectedReapedJobEntries: []*test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Done", LaunchCommand: bgCatCommand1, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received output for echo followed by an entry for the reaped job",
	}

	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test command output not followed by reaped jobs entry
	echoTestCase2 := test_cases.CommandResponseTestCase{
		Command:        fmt.Sprintf("echo %s", echoArgs[1]),
		ExpectedOutput: echoArgs[1],
		SuccessMessage: "✓ Received output for echo",
	}

	if err := echoTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsEmptyTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{},
		SuccessMessage:        "✓ No jobs",
	}

	if err := jobsEmptyTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	sleepCommand := "sleep 10"
	if err := (&test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepCommand,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Output includes job number with PID",
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsOneEntryTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{{
			JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ 1 entry matches the running job",
	}

	if err := jobsOneEntryTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func testBg9Recycle(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	logger := stageHarness.Logger

	fifoPath := filepath.Join(
		"/tmp",
		fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100)),
	)

	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	sleepCommand := "sleep 100"
	sleepCommandTestCase := &test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepCommand,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Output includes job number with PID",
	}
	if err := sleepCommandTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	bgCatCommand := fmt.Sprintf("cat %s", fifoPath)
	if err := (&test_cases.BackgroundCommandResponseTestCase{
		Command:           bgCatCommand,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Output includes job number with PID",
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	if err := WriteToFile(stageHarness, fifoPath, ""); err != nil {
		return err
	}
	time.Sleep(time.Millisecond)

	echoArgument := random.RandomWord()
	echoTestCase := test_cases.CommandResponseWithReapedJobsTestCase{
		Command:               fmt.Sprintf("echo %s", echoArgument),
		ExpectedCommandOutput: echoArgument,
		ExpectedReapedJobEntries: []*test_cases.BackgroundJobStatusEntry{{
			JobNumber: 2, Status: "Done", LaunchCommand: bgCatCommand, Marker: test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ Received output for echo followed by an entry for the reaped job",
	}
	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	sleepCommand2 := "sleep 50"
	if err := (&test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepCommand2,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Output includes job number with PID",
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.PreviousJob},
			{JobNumber: 2, Status: "Running", LaunchCommand: sleepCommand2, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ 2 entries match the running jobs",
	}

	if err := jobsTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
