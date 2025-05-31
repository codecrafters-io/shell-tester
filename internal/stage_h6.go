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

func testH6(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

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

	commandTestCases := []test_cases.CommandResponseTestCase{
		{Command: "echo " + randomWords1, ExpectedOutput: randomWords1, SuccessMessage: commandSuccessMessage},
		{Command: "echo " + randomWords2, ExpectedOutput: randomWords2, SuccessMessage: commandSuccessMessage},
	}

	for _, command := range commandTestCases {
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

	// Test up-arrow navigation (going back in history)
	expectedCommands := []struct {
		command string
		message string
	}{
		{echoCommand.Command, "echo " + randomWords3},
		{randomCommand, randomCommand},
	}

	for _, expected := range expectedCommands {
		if err := shell.SendCommandRaw(upArrow); err != nil {
			return err
		}
		stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<UP ARROW>", expected.message)
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: "$ " + expected.command,
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(`^\s*` + expected.command + `\s*$`),
			},
			StayOnSameLine: true,
		})
		if err := asserter.AssertWithoutPrompt(); err != nil {
			return err
		}
		asserter.PopAssertion()
		stageHarness.Logger.Successf("✓ Prompt line matches %q", expected.message)
	}

	// Down-arrow should go forward to the echo command
	if err := shell.SendCommandRaw(downArrow); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<DOWN ARROW>", echoCommand.Command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + echoCommand.Command,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*` + echoCommand.Command + `\s*$`),
		},
		StayOnSameLine: true,
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	asserter.PopAssertion()
	stageHarness.Logger.Successf("✓ Prompt line matches %q", echoCommand.Command)

	// Execute the echo command
	if err := shell.SendCommandRaw("\n"); err != nil {
		return err
	}
	stageHarness.Logger.Infof("Executing command %q", echoCommand.Command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ " + echoCommand.Command,
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: randomWords3,
	})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}
	stageHarness.Logger.Successf("✓ Command executed with expected output")

	return logAndQuit(asserter, nil)
}
