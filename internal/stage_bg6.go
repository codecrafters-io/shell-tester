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

func testBG6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	fifoBaseNames := random.RandomWords(2)
	fifoPath1 := fmt.Sprintf("/tmp/%s-%d", fifoBaseNames[0], random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath1, 0644); err != nil {
		return err
	}

	fifoPath2 := fmt.Sprintf("/tmp/%s-%d", fifoBaseNames[1], random.RandomInt(1, 100))
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

	// Launch 'grep read pattern' to hang the process indefinitely

	grepPattern1 := random.RandomWord()
	bgGrepCommand1 := fmt.Sprintf("grep -q %s %s", grepPattern1, fifoPath1)
	bgGrepTestCase1 := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand1,
		ExpectedJobNumber: 2,
		SuccessMessage:    "✓ Output includes job number with PID",
	}
	if err := bgGrepTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Launch grep read pattern again
	grepPattern2 := random.RandomWord()
	bgGrepCommand2 := fmt.Sprintf("grep -q %s %s", grepPattern2, fifoPath2)
	bgGrepTestCase2 := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand2,
		ExpectedJobNumber: 3,
		SuccessMessage:    "✓ Output includes job number with PID",
	}

	if err := bgGrepTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	// We write to the first fifo
	if err := WriteToFile(stageHarness, fifoPath1, grepPattern1); err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	// Call jobs for the first time
	jobsBuiltinTestCase1 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.JobsBuiltinOutputEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.UnmarkedJob},
			{JobNumber: 2, Status: "Done", LaunchCommand: bgGrepCommand1, Marker: test_cases.PreviousJob},
			{JobNumber: 3, Status: "Running", LaunchCommand: bgGrepCommand2, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 3 entries in the output",
	}
	if err := jobsBuiltinTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to the second fifo
	if err := WriteToFile(stageHarness, fifoPath2, grepPattern2); err != nil {
		return err
	}
	time.Sleep(time.Millisecond)

	// Call jobs for the second time
	jobsBuiltinTestCase2 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.JobsBuiltinOutputEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.PreviousJob},
			{JobNumber: 3, Status: "Done", LaunchCommand: bgGrepCommand2, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ Received 2 entries in the output",
	}
	if err := jobsBuiltinTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Call jobs again
	jobsBuiltinTestCase3 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.JobsBuiltinOutputEntry{
			{JobNumber: 1, Status: "Running", LaunchCommand: sleepCommand, Marker: test_cases.CurrentJob},
		},
		SuccessMessage: "✓ 1 entry matches the running job",
	}
	if err := jobsBuiltinTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
