package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	invalidCommand := getRandomInvalidCommand()

	// We test a nonexistent command first, just to make sure the logic works in a "loop"
	testCase := test_cases.CommandResponseTestCase{
		Command:          invalidCommand,
		ExpectedOutput:   invalidCommand + ": command not found",
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(`^(bash: )?` + invalidCommand + `: (command )?not found$`)},
		SuccessMessage:   "✓ Received command not found message",
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	refTestCase := test_cases.CommandReflectionTestCase{
		Command:             "exit 0",
		SkipPromptAssertion: true,
	}
	if err := refTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	assertFn := func() error {
		return asserter.AssertionCollection.RunWithPromptAssertion(shell.GetScreenState())
	}
	readErr := shell.ReadUntil(utils.AsBool(assertFn))
	output := virtual_terminal.BuildCleanedRow(shell.GetScreenState()[asserter.GetLastLoggedRowIndex()+1])

	// We're expecting EOF since the program should've terminated
	if !errors.Is(readErr, shell_executable.ErrProgramExited) {
		if readErr == nil {
			return fmt.Errorf("Expected program to exit with 0 exit code, program is still running.")
		} else {
			// TODO: Other than ErrProgramExited, what other errors could we get? Are they user errors or internal errors?
			return fmt.Errorf("Error reading output: %v", readErr)
		}
	}

	isTerminated, exitCode := shell.WaitForTermination()
	if !isTerminated {
		return fmt.Errorf("Expected program to exit with 0 exit code, program is still running.")
	}

	logger.Successf("✓ Program exited successfully")

	if exitCode != 0 {
		return fmt.Errorf("Expected 0 as exit code, got %d", exitCode)
	}

	// Most shells return nothing but bash returns the string "exit" when it exits, we allow both styles
	if len(output) > 0 && strings.TrimSpace(output) != "exit" {
		return fmt.Errorf("Expected no output after exit command, got %q", output)
	}

	logger.Successf("✓ No output after exit command")

	return logAndQuit(asserter, nil)
}
