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
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
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

	// Execute initial commands using CommandResponseTestCase
	randomWords1 := strings.Join(random.RandomWords(2), " ")
	randomWords2 := strings.Join(random.RandomWords(2), " ")
	randomWords3 := strings.Join(random.RandomWords(2), " ")
	randomCommand := getRandomInvalidCommand()

	commands := []test_cases.CommandResponseTestCase{
		{Command: "echo " + randomWords1, ExpectedOutput: randomWords1, SuccessMessage: commandSuccessMessage},
		{Command: "echo " + randomWords2, ExpectedOutput: randomWords2, SuccessMessage: commandSuccessMessage},
	}

	for _, command := range commands {
		if err := command.Run(asserter, shell, stageHarness.Logger); err != nil {
			return err
		}
	}

	// Run invalid command test case
	invalidCommandTestCase := test_cases.InvalidCommandTestCase{
		Command: randomCommand,
	}
	if err := invalidCommandTestCase.Run(asserter, shell, stageHarness.Logger); err != nil {
		return err
	}

	// Run the last echo command
	echoCommand := test_cases.CommandResponseTestCase{
		Command:        "echo " + randomWords3,
		ExpectedOutput: randomWords3,
		SuccessMessage: commandSuccessMessage,
	}
	if err := echoCommand.Run(asserter, shell, stageHarness.Logger); err != nil {
		return err
	}

	// First up-arrow should recall the last command (echo randomWords3)
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<UP ARROW>", "echo "+randomWords3)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo " + randomWords3,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*echo ` + randomWords3 + `\s*$`),
		},
		StayOnSameLine: true,
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()
	stageHarness.Logger.Successf("✓ Prompt line matches %q", "echo "+randomWords3)

	// Second up-arrow should recall the invalid command
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<UP ARROW>", randomCommand)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*` + randomCommand + `\s*$`),
		},
		StayOnSameLine: true,
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()
	stageHarness.Logger.Successf("✓ Prompt line matches %q", randomCommand)

	// Third up-arrow should recall echo randomWords2
	if err := shell.SendCommandRaw(upArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<UP ARROW>", "echo "+randomWords2)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ echo " + randomWords2,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*echo ` + randomWords2 + `\s*$`),
		},
		StayOnSameLine: true,
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()
	stageHarness.Logger.Successf("✓ Prompt line matches %q", "echo "+randomWords2)

	// Down-arrow should go forward to the invalid command
	if err := shell.SendCommandRaw(downArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<DOWN ARROW>", randomCommand)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*` + randomCommand + `\s*$`),
		},
		StayOnSameLine: true,
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()
	stageHarness.Logger.Successf("✓ Prompt line matches %q", randomCommand)

	// Execute the invalid command again
	if err := shell.SendCommandRaw("\n"); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Executing command %q", randomCommand)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + randomCommand,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*` + randomCommand + `\s*$`),
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
	stageHarness.Logger.Successf("✓ Command executed with expected error message")

	return logAndQuit(asserter, nil)
}
