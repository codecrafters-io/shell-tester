// Stage 4: Up-arrow navigation test
// This test checks if the shell supports recalling previous commands using the up arrow key.
// It sends a few commands, then simulates up arrow presses and checks if the correct command is recalled at the prompt.

package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH4(stageHarness *test_case_harness.TestCaseHarness) error {
	// Create a new shell instance and asserter for this test run
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Start the shell and assert that the prompt is shown
	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// 2. Simulate pressing the up arrow key and check if the correct command is recalled
	upArrow := "\x1b[A" // ANSI escape sequence for up arrow
	enter := "\n"
	// Send echo hello
	if err := shell.SendCommand("echo hello"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo hello",
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "hello",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo hello",
	})
	if err := shell.SendCommandRaw(enter); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo hello",
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "hello",
	})

	stageHarness.Logger.Infof("Check for a single command up arrow passed")

	// Checking for 2 up arrow keys

	// Send pwd command
	if err := shell.SendCommand("pwd"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ pwd",
	})

	// Send echo world
	if err := shell.SendCommand("echo world"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})
	if err := shell.SendCommandRaw(enter); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "world",
	})

	stageHarness.Logger.Infof("Check for 2 up arrow keys passed")

	// Checking for 3 up arrow keys

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo world",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ pwd",
	})

	// Send the up arrow key
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo hello",
	})

	// If all assertions pass, finish the test
	if err := shell.SendCommandRaw(enter); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo hello",
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "hello",
	})

	stageHarness.Logger.Infof("Check for 3 up arrow keys passed")

	// Log a success message
	stageHarness.Logger.Successf("✓ Up-arrow navigation works as expected")

	return logAndQuit(asserter, nil)
}
