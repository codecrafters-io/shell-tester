// Stage 4: Up-arrow navigation test
// This test checks if the shell supports recalling previous commands using the up arrow key.
// It sends a few commands, then simulates up arrow presses and checks if the correct command is recalled at the prompt.

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

func testH4(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	upArrow := "\x1b[A"

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

	// Test up-arrow navigation (going back in history)
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

	return logAndQuit(asserter, nil)
}
