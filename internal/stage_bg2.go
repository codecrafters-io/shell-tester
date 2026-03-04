package internal

import (
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBG2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Create a named pipe
	fifoPath := fmt.Sprintf("/tmp/%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if err := CreateRandomFIFOWithTeardown(stageHarness, fifoPath, 0644); err != nil {
		return err
	}

	grepPattern := random.RandomWord()
	bgGrepCommand := fmt.Sprintf("grep %s %s", grepPattern, fifoPath)

	testCase := test_cases.BackgroundCommandResponseTestCase{
		Command:           bgGrepCommand,
		SuccessMessage:    "✓ Received next prompt",
		ExpectedJobNumber: 1,
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Write to the fifo
	if err := WriteToFile(stageHarness, fifoPath, grepPattern); err != nil {
		return err
	}

	// Assert background command output
	backgroundCommandOutputTestCase := test_cases.BackgroundCommandOutputOnlyTestCase{
		ExpectedOutputLines: []string{grepPattern},
		SuccessMessage:      fmt.Sprintf("✓ Output of %s found in the next prompt line", shellescape.Quote(bgGrepCommand)),
	}

	if err := backgroundCommandOutputTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
