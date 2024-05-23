package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
)

func assertShellIsRunning(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := test_cases.NewSilentPromptTestCase("$ ")

	if err := testCase.Run(shell, logger); err != nil {
		return fmt.Errorf("Expected shell to print prompt after last command, but it didn't: %v", err)
	}
	return nil
}
