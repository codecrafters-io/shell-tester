package test_cases

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ExitTestCase is a test case that:
// Sends an exit command to the shell
// Verifies that the shell exits with the expected exit code
// If any error occurs returns the error from the corresponding assertion
type ExitTestCase struct {
	AllowedExitCodes []int

	// ShouldSkipSuccessMessage determines if the success message should be skipped (not used just yet, but can be used in the future)
	ShouldSkipSuccessMessage bool
}

func (t ExitTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// First run a command reflection test to verify the command is sent correctly
	commandTestCase := CommandWithNoResponseTestCase{
		Command:             "exit",
		SkipPromptAssertion: true,
	}
	if err := commandTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	assertFn := func() error {
		return asserter.AssertionCollection.RunWithPromptAssertion(shell.GetScreenState())
	}
	readErr := shell.ReadUntilConditionOrTimeout(utils.AsBool(assertFn), logged_shell_asserter.SUBSEQUENT_READ_TIMEOUT)
	output := shell.GetScreenState().GetRow(asserter.GetLastLoggedRowIndex() + 1).String()

	asserter.LogRemainingOutput()

	// We're expecting EOF since the program should've terminated
	if !errors.Is(readErr, shell_executable.ErrProgramExited) {
		if readErr == nil {
			return fmt.Errorf("Expected program to exit, program is still running.")
		} else if errors.Is(readErr, condition_reader.ErrConditionNotMet) {
			return fmt.Errorf("Expected program to exit, program is still running.")
		} else {
			return fmt.Errorf("Error reading output: %v", readErr)
		}
	}

	isTerminated, exitCode := shell.WaitForTermination()
	if !isTerminated {
		return fmt.Errorf("Expected program to exit, program is still running.")
	}

	// We want to be lenient since:
	// - calling `exit` without arguments returns the exit status of the last executed command,
	// - but we don't want to burden users with this requirement.
	if len(t.AllowedExitCodes) == 0 {
		t.AllowedExitCodes = []int{0}
	}
	if !slices.Contains(t.AllowedExitCodes, exitCode) {
		return fmt.Errorf("Expected exit code to be one of %v, got %d", t.AllowedExitCodes, exitCode)
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
