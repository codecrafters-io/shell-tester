package assertions

import (
	"fmt"
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
	s.Logger.Debugf("--------------------------------")
	for _, row := range s.Shell.GetScreenState() {
		cleanedRow := buildCleanedRow(row)
		if len(cleanedRow) > 0 {
			s.Logger.Debugf(cleanedRow)
		}
	}
	s.Logger.Debugf("--------------------------------")
}

func (s ScreenAsserter) LogCurrentRow() {
	s.Logger.Debugf("--------------------------------")
	cleanedRow := buildCleanedRow(s.Shell.GetScreenState()[s.rowIndex])
	if len(cleanedRow) > 0 {
		s.Logger.Debugf(cleanedRow)
	}
	s.Logger.Debugf("--------------------------------")
}

func (s *ScreenAsserter) LogUptoCurrentRow() {
	for i := s.loggedUptoRowIndex; i <= s.rowIndex; i++ {
		s.LogRow(i)
	}
	s.UpdateLoggedUptoRowIndex()
}

func (s *ScreenAsserter) LogRow(rowIndex int) {
	s.Logger.Debugf("--------------------------------")
	cleanedRow := buildCleanedRow(s.Shell.GetScreenState()[rowIndex])
	if len(cleanedRow) > 0 {
		s.Logger.Debugf(cleanedRow)
	}
	s.Logger.Debugf("--------------------------------")
}

func (s *ScreenAsserter) UpdateLoggedUptoRowIndex() {
	s.loggedUptoRowIndex = s.rowIndex
}

func (s ScreenAsserter) PromptAssertion(rowIndex int, expectedPrompt string) PromptAssertion {
	return PromptAssertion{rowIndex: rowIndex, expectedPrompt: expectedPrompt, screenAsserter: &s}
}

func (s ScreenAsserter) SingleLineAssertion(rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) SingleLineScreenStateAssertion {
	return SingleLineScreenStateAssertion{rowIndex: rowIndex, expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation, screenAsserter: &s}
}

func (s *ScreenAsserter) AddAssertion(assertion Assertion) {
	s.Assertions = append(s.Assertions, assertion)
}

func (s *ScreenAsserter) UpdateRowIndex(increment int) {
	fmt.Println("Updating row index", increment)
	s.rowIndex += increment
	fmt.Println(s.rowIndex)
}

func (s *ScreenAsserter) RunAllAssertions() error {
	for _, assertion := range s.Assertions {
		if err := assertion.Run(); err != nil {
			return err
		}
		assertion.UpdateRowIndex()
		fmt.Println(s.rowIndex)
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
