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

func testBG5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Launch a grep fifo
	fifoPath := fmt.Sprintf("/tmp/%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath, 0644); err != nil {
		return err
	}

	grepPattern := random.RandomWord()
	bgGrepCommand := fmt.Sprintf("grep -q %s %s", grepPattern, fifoPath)

	bgJobTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand,
		ExpectedJobNumber: 1,
		SuccessMessage:    fmt.Sprintf("✓ Expected entry found for the started job"),
	}

	if err := bgJobTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Call jobs
	jobsBuiltinTestCase1 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.JobsBuiltinOutputEntry{{
			JobNumber:     1,
			Status:        "Running",
			LaunchCommand: bgGrepCommand,
			Marker:        test_cases.CurrentJob,
		}},
		SuccessMessage: "Expected 1 entry found for the running job",
	}

	if err := jobsBuiltinTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to fifo
	if err := WriteToFile(stageHarness, fifoPath, grepPattern); err != nil {
		return err
	}

	// A small delay since grep takes some time to process and exit
	time.Sleep(time.Millisecond)

	// Call jobs again
	jobsBuiltinTestCase2 := test_cases.JobsBuiltinResponseTestCase{
		ExpectedOutputEntries: []test_cases.JobsBuiltinOutputEntry{{
			JobNumber:     1,
			Status:        "Done",
			LaunchCommand: bgGrepCommand,
			Marker:        test_cases.CurrentJob,
		}},
		SuccessMessage: "Expected 1 entry found for the finished job",
	}

	if err := jobsBuiltinTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
