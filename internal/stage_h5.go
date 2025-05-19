// Stage 5: Down-arrow navigation test
// This test checks if the shell supports traversing through history with both up and down arrow keys.
// It sends a few commands, simulates up and down arrow presses, and checks if the correct command is recalled at the prompt.

package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH5(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	upArrow := "\x1b[A"
	downArrow := "\x1b[B"

	// $ ls dist/
	if err := shell.SendCommand("ls dist/"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls dist/"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// $ cd dist/
	if err := shell.SendCommand("cd dist/"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ cd dist/"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// $ ls
	if err := shell.SendCommand("ls"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// <UP ARROW> (should recall 'ls')
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	// <UP ARROW> (should recall 'cd dist/')
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	// <DOWN ARROW> (should go forward to 'ls')
	if err := shell.SendCommandRaw(downArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<DOWN ARROW>")
	// <ENTER> (should execute 'ls' again)
	if err := shell.SendCommandRaw("\n"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
