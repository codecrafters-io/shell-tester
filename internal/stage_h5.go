// Stage 5: Down-arrow navigation test
// This test checks if the shell supports traversing through history with both up and down arrow keys.
// It sends a few commands, simulates up and down arrow presses, and checks if the correct command is recalled at the prompt.

package internal

import (
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/random"
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

	// $ echo hello
	randomWords1 := strings.Join(random.RandomWords(2), " ")
	if err := shell.SendCommand("echo " + randomWords1); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ echo " + randomWords1})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: randomWords1})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// $ echo world
	randomWords2 := strings.Join(random.RandomWords(2), " ")
	if err := shell.SendCommand("echo " + randomWords2); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ echo " + randomWords2})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: randomWords2})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// $ nonexistent_command
	randomCommand := getRandomInvalidCommand()
	if err := shell.SendCommand(randomCommand); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ " + randomCommand})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "command not found: " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^.*command not found.*` + randomCommand + `.*$`),
			regexp.MustCompile(`^.*` + randomCommand + `.*command not found.*$`),
			regexp.MustCompile(`^.*` + randomCommand + `.*not found.*$`),
			regexp.MustCompile(`^(zsh|bash): command not found: ` + randomCommand + `$`),
		},
	})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// $ echo random words
	randomWords3 := strings.Join(random.RandomWords(2), " ")
	if err := shell.SendCommand("echo " + randomWords3); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: "$ echo " + randomWords3})
	asserter.AddAssertion(assertions.SingleLineAssertion{ExpectedOutput: randomWords3})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// <UP ARROW> (should recall 'echo ' + randomWords3)
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	// <UP ARROW> (should recall 'nonexistent_command')
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	// <UP ARROW> (should recall 'echo ' + randomWords2)
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<UP ARROW>")
	// <DOWN ARROW> (should go forward to 'nonexistent_command')
	if err := shell.SendCommandRaw(downArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("<DOWN ARROW>")
	// <ENTER> (should execute 'nonexistent_command' again)
	if err := shell.SendCommandRaw("\n"); err != nil {
		return err
	}
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*` + randomCommand + `$`),
		},
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "command not found: " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^.*command not found.*` + randomCommand + `.*$`),
			regexp.MustCompile(`^.*` + randomCommand + `.*command not found.*$`),
			regexp.MustCompile(`^.*` + randomCommand + `.*not found.*$`),
			regexp.MustCompile(`^(zsh|bash): command not found: ` + randomCommand + `$`),
		},
	})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
