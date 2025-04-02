package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandMultipleCompletionsTestCase is a test case that:
// Sends a command to the shell
// Asserts that the prompt line reflects the command
// Sends TAB
// Asserts that the expected reflection is printed to the screen (with a space after it)
// If any error occurs returns the error from the corresponding assertion
type CommandMultipleCompletionsTestCase struct {
	// RawCommand is the command to send to the shell
	RawCommand string

	// ExpectedReflection is the custom reflection to use
	ExpectedReflection string

	// ExpectedAutocompletedReflectionHasNoSpace is true if
	// the expected reflection should have no space after it
	ExpectedAutocompletedReflectionHasNoSpace bool

	// CheckForBell is true if we should check for a bell
	CheckForBell bool

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// tabCount is the number of tabs to send after the command
	TabCount int

	// SkipPromptAssertion is a flag to skip the final prompt assertion
	SkipPromptAssertion bool
}

func (t CommandMultipleCompletionsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Log the details of the command before sending it
	logCommand(logger, t.RawCommand)

	// Send the command to the shell
	if err := shell.SendCommandRaw(t.RawCommand); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	render(shell)

	inputReflection := fmt.Sprintf("$ %s", t.RawCommand)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: false,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	render(shell)

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", inputReflection)

	// // The space at the end of the reflection won't be present, so replace that assertion
	// asserter.PopAssertion()

	// Send TAB
	for i := 0; i < t.TabCount; i++ {
		shouldRingBell := (i == 0 && t.CheckForBell)
		logTab(logger, t.ExpectedReflection, shouldRingBell)

		time.Sleep(1 * time.Millisecond)

		if err := shell.SendCommandRaw("\t"); err != nil {
			return fmt.Errorf("Error sending command to shell: %v", err)
		}
	}

	render(shell)

	if t.CheckForBell {
		bellChannel := shell.VTBellChannel()
		asserter.AddAssertion(assertions.BellAssertion{
			BellChannel: bellChannel,
		})
		// Run the assertion, before sending the enter key
		if err := asserter.AssertWithoutPrompt(); err != nil {
			return err
		}

		logger.Successf("✓ Received bell")
		// Pop the bell assertion after running
		asserter.PopAssertion()
	}

	render(shell)

	commandReflection := t.ExpectedReflection
	// Space after autocomplete
	if !t.ExpectedAutocompletedReflectionHasNoSpace {
		commandReflection = fmt.Sprintf("%s ", t.ExpectedReflection)
	}

	render(shell)

	// Assert auto-completion
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	render(shell)

	// Only if we attempted to autocomplete, print the success message
	logger.Successf("✓ Prompt line matches %q", t.ExpectedReflection)

	// The space at the end of the reflection won't be present, so replace that assertion
	// asserter.PopAssertion()

	render(shell)

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: inputReflection,
		StayOnSameLine: true,
	})
	// Run the assertion, before sending the enter key
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	render(shell)

	var assertFuncToRun func() error
	if t.SkipPromptAssertion {
		assertFuncToRun = asserter.AssertWithoutPrompt
	} else {
		assertFuncToRun = asserter.AssertWithPrompt
	}

	render(shell)

	if err := assertFuncToRun(); err != nil {
		return err
	}

	render(shell)

	logger.Successf("%s", t.SuccessMessage)
	return nil
}

func render(shell *shell_executable.ShellExecutable) {
	screen := (shell.GetScreenState())
	for row := range screen {
		cleanedRow := virtual_terminal.BuildCleanedRow2(screen[row])
		if len(cleanedRow) > 0 {
			fmt.Println(cleanedRow)
		}
	}
}
