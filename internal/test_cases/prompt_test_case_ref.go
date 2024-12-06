package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// TODO: Remove PromptTestCase entirely, replace with PromptAssertion invoked within ScreenAsserter
// PromptTestCase verifies a prompt exists, and that there's no extra output after it.
type PromptTestCase struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	// shouldOmitSuccessLog determines whether a success log should be written.
	// When re-using this test case within other higher-order test cases, emitting success logs
	// all the time can get pretty noisy.
	shouldOmitSuccessLog bool
}

func NewPromptTestCase(expectedPrompt string) PromptTestCase {
	return PromptTestCase{expectedPrompt: expectedPrompt}
}

func NewSilentPromptTestCase(expectedPrompt string) PromptTestCase {
	return PromptTestCase{expectedPrompt: expectedPrompt, shouldOmitSuccessLog: true}
}

// 1. PromptAssertion
//
// $
//
// ----
//
// writeText("echo hello")
// 1. CommandReflectionAssertion
// 2. SingleLineOutputAssertion
// 3. PromptAssertion
func (t PromptTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	return nil
}
