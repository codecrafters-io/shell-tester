package test_cases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ExitTestCase is a test case that:
// Sends an exit command to the shell
// Verifies that the shell exits with the expected exit code
// If any error occurs returns the error from the corresponding assertion
type ExitTestCase struct {
	// Command is the exit command to send to the shell (e.g. "exit 0")
	Command string

	// ExpectedExitCode is the expected exit code
	ExpectedExitCode int

	// ShouldSkipSuccessMessage determines if the success message should be skipped (not used just yet, but can be used in the future)
	ShouldSkipSuccessMessage bool
}

func (t ExitTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// First run a command reflection test to verify the command is sent correctly
	commandTestCase := CommandWithNoResponseTestCase{
		Command:             t.Command,
		SkipPromptAssertion: true,
	}
	if err := commandTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	assertFn := func() error {
		return asserter.AssertionCollection.RunWithPromptAssertion(shell.GetScreenState())
	}
	readErr := shell.ReadUntilConditionOrTimeout(utils.AsBool(assertFn), logged_shell_asserter.SUBSEQUENT_READ_TIMEOUT)
	output := virtual_terminal.BuildCleanedRow(shell.GetScreenState()[asserter.GetLastLoggedRowIndex()+1])

	asserter.LogRemainingOutput()

	// We're expecting EOF since the program should've terminated
	if !errors.Is(readErr, shell_executable.ErrProgramExited) {
		if readErr == nil {
			return fmt.Errorf("Expected program to exit with %d exit code, program is still running.", t.ExpectedExitCode)
		} else if errors.Is(readErr, condition_reader.ErrConditionNotMet) {
			return fmt.Errorf("Expected program to exit with %d exit code, program is still running.", t.ExpectedExitCode)
		} else {
			return fmt.Errorf("Error reading output: %v", readErr)
		}
	}

	isTerminated, exitCode := shell.WaitForTermination()
	if !isTerminated {
		return fmt.Errorf("Expected program to exit with %d exit code, program is still running.", t.ExpectedExitCode)
	}

	if exitCode != t.ExpectedExitCode {
		return fmt.Errorf("Expected %d as exit code, got %d", t.ExpectedExitCode, exitCode)
	}

	// Most shells return nothing but bash returns the string "exit" when it exits, we allow both styles
	if len(output) > 0 && strings.TrimSpace(output) != "exit" {
		// If there is some unexpected output, we need to log it before returning an error
		asserter.LogRemainingOutput()
		return fmt.Errorf("Expected no output after exit command, got %q", output)
	}

	if !t.ShouldSkipSuccessMessage {
		logger.Successf("✓ Program exited successfully")
	}
	return nil
}
