package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// ToDo: This is a prototype, think about edge cases + implement prompt test case specifically
// TODO: Remove CommandResponseTestCase entirely, replace with SingleLineOutputAssertion invoked within ScreenAsserter
// CommandResponseTestCase reads the output from the shell, and verifies that it matches the expected output.
type CommandResponseTestCase struct {
	// command is the command that will be sent to the shell
	command string
}

func NewCommandResponseTestCase(command string) CommandResponseTestCase {
	return CommandResponseTestCase{command: command}
}

func (t CommandResponseTestCase) Run(screenAsserter *assertions.ScreenAsserter, shouldOmitSuccessLog bool) error {
	err := screenAsserter.Shell.SendCommand(t.command)
	if err != nil {
		return fmt.Errorf("Error sending command: %v", err)
	}

	err = screenAsserter.Shell.ReadUntil(screenAsserter.RunBool)

	screenAsserter.RunBool()

	if err != nil {
		// If the user sent any output, let's print it before the error message.
		if len(screenAsserter.Shell.GetScreenState()) > 0 {
			screenAsserter.LogFullScreenState()
		}

		return fmt.Errorf("Expected prompt (%q) to be printed, got %q", t.command, utils.BuildCleanedRow(screenAsserter.Shell.GetScreenState()[0]))
	}

	err = screenAsserter.Shell.ReadUntilTimeout(10 * time.Millisecond)

	// Whether the value matches our expectations or not, we print it
	// screenAsserter.LogUptoCurrentRow()

	// We failed to read extra output
	if err != nil {
		return fmt.Errorf("Error reading output: %v", err)
	}

	if !shouldOmitSuccessLog {
		screenAsserter.Logger.Successf("✓ Received prompt")
	}

	return nil
}
