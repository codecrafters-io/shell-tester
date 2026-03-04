package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"al.essio.dev/pkg/shellescape"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/codecrafters-io/tester-utils/testing"
)

func testBG3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Create a named pipe
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

	bgCatCommand := fmt.Sprintf("cat %s", fifoPath1)
	bgCommandTestCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgCatCommand,
		SuccessMessage:    "✓ Received next prompt",
		ExpectedJobNumber: 1,
	}

	if err := bgCommandTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Launch a foreground test case that reads from the fifo
	fgCatCommand := fmt.Sprintf("cat %s", fifoPath2)
	fgCommandTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:             fgCatCommand,
		SkipPromptAssertion: true,
	}

	if err := fgCommandTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	// Write to the fifo 1, and assert the output
	fifo1Contents := "Hello from FIFO #1\n"
	if err := WriteToFile(stageHarness, fifoPath1, fifo1Contents); err != nil {
		return err
	}

	// Assert background cat command output
	bgCommandOutputTestCasae := test_cases.OutputOnlyTestCase{
		ExpectedOutputLines: []string{strings.TrimSuffix(fifo1Contents, "\n")},
		SuccessMessage:      fmt.Sprintf("✓ Received output from background job %s", shellescape.Quote(bgCatCommand)),
	}

	if err := bgCommandOutputTestCasae.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to the fifo 2, and assert the output
	fifo2Contents := "Hello from FIFO #2\n"
	if err := WriteToFile(stageHarness, fifoPath2, fifo2Contents); err != nil {
		return err
	}

	// Assert foreground cat command output
	fgCommandOutputTestCase := test_cases.OutputOnlyTestCase{
		ExpectedOutputLines: []string{strings.TrimSuffix(fifo2Contents, "\n")},
		SuccessMessage:      fmt.Sprintf("✓ Received output from foreground job %s", shellescape.Quote(fgCatCommand)),
	}

	if err := fgCommandOutputTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Sleep before logging so that the finished job status and prompt are guaranteed to be printed in the fixtures
	if testing.IsRecordingOrEvaluatingFixtures() {
		time.Sleep(time.Second)
	}

	return logAndQuit(asserter, nil)
}
