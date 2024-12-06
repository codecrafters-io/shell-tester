package assertions

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ToDo: Prototype, not yet sure what the ideal consituents of this struct should be
type ScreenAsserter struct {
	Shell  *shell_executable.ShellExecutable
	Logger *logger.Logger
}

func NewScreenAsserter(shell *shell_executable.ShellExecutable, logger *logger.Logger) ScreenAsserter {
	return ScreenAsserter{Shell: shell, Logger: logger}
}

func (s ScreenAsserter) LogFullScreenState() {
	for _, row := range s.Shell.GetScreenState() {
		s.Logger.Debugf(strings.Join(row, ""))
	}
}
