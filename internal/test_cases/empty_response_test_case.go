package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// EmptyResponseTestCase verifies a prompt exists, sends a command and verifies that there is no output.
type EmptyResponseTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// SuccessMessage is logged if the response is empty
	SuccessMessage string
}

func (t EmptyResponseTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger, skipReadingPrompt bool) error {
	promptTestCase := NewSilentPromptTestCase("$ ")

	if err := shell.SendCommand(t.Command); err != nil {
		return err
	}

	if !skipReadingPrompt {
		if err := promptTestCase.Run(shell, logger); err != nil {
			return err
		}

		logger.Successf("✓ %s", t.SuccessMessage)
	}
	return nil
}
