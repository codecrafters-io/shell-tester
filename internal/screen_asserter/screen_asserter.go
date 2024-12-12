package screen_asserter

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/logger"
)

// Logging success rows
// Reading from shell
// Running assertions
type ScreenAsserter struct {
	Shell      *shell_executable.ShellExecutable
	Logger     *logger.Logger
	Assertions []assertions.Assertion

	nextRowToLog int
}

func NewScreenAsserter(shell *shell_executable.ShellExecutable, logger *logger.Logger) *ScreenAsserter {
	return &ScreenAsserter{Shell: shell, Logger: logger}
}

func (s *ScreenAsserter) RunWithPromptAssertion() error {
	return s.RunWithExtraAssertions([]assertions.Assertion{
		assertions.PromptAssertion{ExpectedPrompt: "$ "},
	})
}

func (s *ScreenAsserter) Run() error {
	return s.RunWithExtraAssertions(nil)
}

func (s *ScreenAsserter) RunWithExtraAssertions(extraAssertions []assertions.Assertion) error {
	allAssertions := append(s.Assertions, extraAssertions...)
	currentRowIndex := 0

	for _, assertion := range allAssertions {
		processedRowCount, err := assertion.Run(s.Shell.GetScreenState(), currentRowIndex)
		if err != nil {
			return err
		}

		currentRowIndex += processedRowCount

		// Log "success" rows that were processed
		if s.nextRowToLog <= currentRowIndex {
			for i := s.nextRowToLog; i < currentRowIndex; i++ {
				s.logRow(i)
			}
			s.nextRowToLog = currentRowIndex
		}
	}

	return nil
}

func (s *ScreenAsserter) PushAssertion(assertion assertions.Assertion) {
	s.Assertions = append(s.Assertions, assertion)
}

func (s *ScreenAsserter) logRow(rowIndex int) {
	cleanedRow := utils.BuildCleanedRow(s.Shell.GetScreenState()[rowIndex])
	if len(cleanedRow) > 0 {
		// s.Logger.Debugf(cleanedRow)
		s.Shell.LogOutput([]byte(cleanedRow))
	} else {
		// ToDo: Remove this, this is an assertion for rowIndex
		// values not going out of range
		s.Logger.Debugf("No output")
	}
}
