// Stage 4: Up-arrow navigation test
// This test checks if the shell supports recalling previous commands using the up arrow key.
// It sends a few commands, then simulates up arrow presses and checks if the correct command is recalled at the prompt.

package internal

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testH4(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	upArrow := "\x1b[A"

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
		output  string
	}{
		{echoCommand.Command, "echo " + randomWords3, randomWords3},
		{randomCommand, randomCommand, ""},
		{"echo " + randomWords2, "echo " + randomWords2, randomWords2},
	}

	for _, expected := range expectedCommands {
		if err := shell.SendCommandRaw(upArrow); err != nil {
			return err
		}
		stageHarness.Logger.Infof("Pressed %q (expecting to recall %q)", "<UP ARROW>", expected.message)
		asserter.AddAssertion(&assertions.SingleLineAssertion{
			ExpectedOutput: "$ " + expected.command,
			StayOnSameLine: true,
		})
		if err := asserter.AssertWithoutPrompt(); err != nil {
			return err
		}
		asserter.PopAssertion()
		stageHarness.Logger.Successf("✓ Prompt line matches %q", expected.message)
	}

	return logAndQuit(asserter, nil)
}
