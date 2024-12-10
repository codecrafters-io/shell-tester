package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
)

// ToDo: This is a prototype, think about edge cases + implement prompt test case specifically
// TODO: Remove ResponseTestCase entirely, replace with SingleLineOutputAssertion invoked within ScreenAsserter
// ResponseTestCase reads the output from the shell, and verifies that it matches the expected output.
type ResponseTestCase struct {
}

func NewResponseTestCase() ResponseTestCase {
	return ResponseTestCase{}
}

func (t ResponseTestCase) Run(screenAsserter *assertions.ScreenAsserter, shouldOmitSuccessLog bool) error {
	err := screenAsserter.Shell.ReadUntil(screenAsserter.WrappedRunAllAssertions)
	// If assertions contain a single assertion and if that is a prompt assertion, we need to log current row else pass

	screenAsserter.RunAllAssertions(false)

	if err != nil {
		// If the user sent any output, let's print it before the error message.
		if len(screenAsserter.Shell.GetScreenState()) > 0 {
			screenAsserter.LogFullScreenState()
		}

		// ToDo: Figure out how to get the expected output here
		return fmt.Errorf("Expected response %q, got %q", "", buildCleanedRow(screenAsserter.Shell.GetScreenState()[0]))
	}

	err = screenAsserter.Shell.ReadUntilTimeout(10 * time.Millisecond)

	// Whether the value matches our expectations or not, we print it
	// fmt.Println("Before logging in response test case", screenAsserter.GetRowIndex(), screenAsserter.GetLoggedUptoRowIndex())
	screenAsserter.LogUptoCurrentRow()
	// fmt.Println("After logging in response test case", screenAsserter.GetRowIndex(), screenAsserter.GetLoggedUptoRowIndex())

	// We failed to read extra output
	if err != nil {
		return fmt.Errorf("Error reading output: %v", err)
	}

	if !shouldOmitSuccessLog {
		screenAsserter.Logger.Successf("✓ Received prompt")
	}

	return nil
}
