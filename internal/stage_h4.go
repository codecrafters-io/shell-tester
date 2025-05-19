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
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	upArrow := "\x1b[A"

	// $ ls dist/
	if err := shell.SendCommand("ls dist/"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls dist/"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// Send cd dist/
	if err := shell.SendCommand("cd dist/"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ cd dist/"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// Send ls
	if err := shell.SendCommand("ls"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// Send up arrow
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls"})

	// Send up arrow again
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls dist/"})

	// Send enter
	if err := shell.SendCommandRaw("\n"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ ls dist/"})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "main.out"})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
