package internal

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBG8(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	if err := testBg8ResetToZero(stageHarness); err != nil {
		return err
	}

	logger.Infof("Tearing down shell")

	return testBg8Recycle(stageHarness)
}

func testBg8ResetToZero(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	logger := stageHarness.Logger

	fifoPath1 := fmt.Sprintf("/tmp/%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath1, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	grepPattern1 := random.RandomWord()
	bgGrepCommand1 := fmt.Sprintf("grep -q %s %s", grepPattern1, fifoPath1)
	bgGrepTestCase := &test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand1,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Received entry for the launched job",
	}
	if err := bgGrepTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to fifo to make job 1 'Done'
	if err := WriteToFile(stageHarness, fifoPath1, grepPattern1); err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	echoArg := random.RandomWord()
	echoTestCase := test_cases.CommandResponseWithReapedJobsTestCase{
		Command:               fmt.Sprintf("echo %s", echoArg),
		ExpectedCommandOutput: echoArg,
		ExpectedReapedJobEntries: []*test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Done", LaunchCommand: bgGrepCommand1, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received command output followed by an entry for the reaped job",
	}

	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsEmptyTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{},
		SuccessMessage:        "✓ Received no entries",
	}

	if err := jobsEmptyTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	sleepCommand := "sleep 10"
	if err := (&test_cases.BackgroundCommandResponseTestCase{Command: sleepCommand, ExpectedJobNumber: 1, SuccessMessage: "✓ Received entry for the launched job"}).Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsOneEntryTestCase := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{{
			JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ Received 1 entry for the running job",
	}

	if err := jobsOneEntryTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func testBg8Recycle(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	logger := stageHarness.Logger

	fifoPath := fmt.Sprintf("/tmp/%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	sleepLong := "sleep 100"
	bgLongCmd := &test_cases.BackgroundCommandResponseTestCase{
		Command:           sleepLong,
		ExpectedJobNumber: 1,
		SuccessMessage:    "✓ Received entry for the launched job",
	}
	if err := bgLongCmd.Run(asserter, shell, logger); err != nil {
		return err
	}

	grepPattern := random.RandomWord()
	bgGrepCommand := fmt.Sprintf("grep -q %s %s", grepPattern, fifoPath)
	if err := (&test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Received entry for the launched job",
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	if err := WriteToFile(stageHarness, fifoPath, grepPattern); err != nil {
		return err
	}
	time.Sleep(time.Millisecond)

	echoArgument := random.RandomWord()
	echoTestCase := test_cases.CommandResponseWithReapedJobsTestCase{
		Command:               fmt.Sprintf("echo %s", echoArgument),
		ExpectedCommandOutput: echoArgument,
		ExpectedReapedJobEntries: []*test_cases.BackgroundJobStatusEntry{{
			JobNumber: 2, Status: "Done", LaunchCommand: bgGrepCommand, Marker: test_cases.CurrentJob,
		}},
		SuccessMessage: "✓ Received command output followed by an entry for the reaped job",
	}
	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	sleep50 := "sleep 50"
	if err := (&test_cases.BackgroundCommandResponseTestCase{
		Command:           sleep50,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Received entry for the launched job",
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	jobsTwo := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepLong, Marker: test_cases.PreviousJob},
			{JobNumber: 2, Status: "Running", LaunchCommand: sleep50, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 2 entries for the running jobs",
	}

	if err := jobsTwo.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
