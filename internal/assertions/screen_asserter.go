package assertions

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ToDo: Prototype, not yet sure what the ideal consituents of this struct should be
type ScreenAsserter struct {
	Shell      *shell_executable.ShellExecutable
	Logger     *logger.Logger
	Assertions []Assertion
	// This is the cursor we will use for selecting the row to assert on ideally
	// For now this is only used for logging
	rowIndex           int
	loggedUptoRowIndex int
}

func NewScreenAsserter(shell *shell_executable.ShellExecutable, logger *logger.Logger) *ScreenAsserter {
	return &ScreenAsserter{Shell: shell, Logger: logger, rowIndex: 0, loggedUptoRowIndex: 0}
}

func (s ScreenAsserter) LogFullScreenState() {
	for _, row := range s.Shell.GetScreenState() {
		cleanedRow := buildCleanedRow(row)
		if len(cleanedRow) > 0 {
			s.Logger.Debugf(cleanedRow)
		}
	}
}

func (s ScreenAsserter) LogCurrentRow() {
	cleanedRow := buildCleanedRow(s.Shell.GetScreenState()[s.rowIndex])
	if len(cleanedRow) > 0 {
		s.Logger.Debugf(cleanedRow)
	}
}

func (s *ScreenAsserter) LogUptoCurrentRow() {
	if s.loggedUptoRowIndex == s.rowIndex {
		return
	}
	for i := s.loggedUptoRowIndex; i <= s.rowIndex; i++ {
		s.LogRow(i)
	}
	s.UpdateLoggedUptoRowIndex()
}

func (s *ScreenAsserter) LogRow(rowIndex int) {
	cleanedRow := buildCleanedRow(s.Shell.GetScreenState()[rowIndex])
	if len(cleanedRow) > 0 {
		s.Logger.Debugf(cleanedRow)
	}
}

func (s *ScreenAsserter) UpdateLoggedUptoRowIndex() {
	s.loggedUptoRowIndex = s.rowIndex + 1
}

func (s ScreenAsserter) PromptAssertion(rowIndex int, expectedPrompt string, screenAsserter *ScreenAsserter) PromptAssertion {
	return PromptAssertion{rowIndex: rowIndex, expectedPrompt: expectedPrompt, screenAsserter: screenAsserter}
}

func (s ScreenAsserter) SingleLineAssertion(rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string, screenAsserter *ScreenAsserter) SingleLineScreenStateAssertion {
	return SingleLineScreenStateAssertion{rowIndex: rowIndex, expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation, screenAsserter: screenAsserter}
}

func (s *ScreenAsserter) AddAssertion(assertion Assertion) {
	s.Assertions = append(s.Assertions, assertion)
}

func (s *ScreenAsserter) UpdateRowIndex(increment int) {
	s.rowIndex += increment
}

func (s *ScreenAsserter) RunAllAssertions() error {
	for _, assertion := range s.Assertions {
		if err := assertion.Run(); err != nil {
			return err
		}
		assertion.UpdateRowIndex()
	}
	return nil
}

func (s *ScreenAsserter) WrappedRunAllAssertions() bool {
	// True if the prompt assertion is a success
	return s.RunAllAssertions() == nil
}

func (s *ScreenAsserter) ClearAssertions() {
	s.Assertions = []Assertion{}
}

func (s *ScreenAsserter) GetRowIndex() int {
	return s.rowIndex
}

func (s *ScreenAsserter) GetLoggedUptoRowIndex() int {
	return s.loggedUptoRowIndex
}
