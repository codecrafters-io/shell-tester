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
}

func NewScreenAsserter(shell *shell_executable.ShellExecutable, logger *logger.Logger) *ScreenAsserter {
	return &ScreenAsserter{Shell: shell, Logger: logger}
}

func (s ScreenAsserter) LogFullScreenState() {
	for _, row := range s.Shell.GetScreenState() {
		cleanedRow := buildCleanedRow(row)
		if len(cleanedRow) > 0 {
			s.Logger.Debugf(cleanedRow)
		}
	}
}

func (s ScreenAsserter) PromptAssertion(rowIndex int, expectedPrompt string, shouldOmitSuccessLog bool) PromptAssertion {
	return PromptAssertion{rowIndex: rowIndex, expectedPrompt: expectedPrompt, screenAsserter: &s, shouldOmitSuccessLog: shouldOmitSuccessLog}
}

func (s ScreenAsserter) SingleLineAssertion(rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) SingleLineScreenStateAssertion {
	return SingleLineScreenStateAssertion{rowIndex: rowIndex, expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation, screenAsserter: &s}
}

func (s *ScreenAsserter) AddAssertion(assertion Assertion) {
	s.Assertions = append(s.Assertions, assertion)
}

func (s *ScreenAsserter) RunAllAssertions() error {
	for _, assertion := range s.Assertions {
		if err := assertion.Run(); err != nil {
			return err
		}
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
