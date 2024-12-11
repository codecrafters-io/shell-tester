package assertions

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/logger"
)

// ToDo: Prototype, not yet sure what the ideal consituents of this struct should be
type ScreenAsserter struct {
	Shell      *shell_executable.ShellExecutable
	Logger     *logger.Logger
	Assertions []Assertion

	nextRowToLog int
}

func NewScreenAsserter(shell *shell_executable.ShellExecutable, logger *logger.Logger) *ScreenAsserter {
	return &ScreenAsserter{Shell: shell, Logger: logger}
}

func (s *ScreenAsserter) LogFullScreenState() {
	for _, row := range s.Shell.GetScreenState() {
		cleanedRow := utils.BuildCleanedRow(row)
		if len(cleanedRow) > 0 {
			s.Logger.Debugf(cleanedRow)
		}
	}
}

// func (s *ScreenAsserter) LogCurrentRow() {
// 	cleanedRow := utils.BuildCleanedRow(s.Shell.GetScreenState()[s.rowIndex])
// 	if len(cleanedRow) > 0 {
// 		s.Logger.Debugf(cleanedRow)
// 	}
// }

// func (s *ScreenAsserter) LogUptoCurrentRow() {
// 	for i := s.loggedUptoRowIndex; i <= s.rowIndex; i++ {
// 		s.LogRow(i)
// 	}
// 	s.UpdateLoggedUptoRowIndex()
// }

func (s *ScreenAsserter) LogRow(rowIndex int) {
	cleanedRow := utils.BuildCleanedRow(s.Shell.GetScreenState()[rowIndex])
	if len(cleanedRow) > 0 {
		s.Logger.Debugf(cleanedRow)
	} else {
		// ToDo: Remove this, this is an assertion for rowIndex
		// values not going out of range
		s.Logger.Debugf("No output")
	}
}

// func (s *ScreenAsserter) UpdateLoggedUptoRowIndex() {
// 	s.loggedUptoRowIndex = s.rowIndex + 1
// }

func (s *ScreenAsserter) PromptAssertion(expectedPrompt string) PromptAssertion {
	return NewPromptAssertion(expectedPrompt)
}

func (s *ScreenAsserter) SingleLineAssertion(rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) SingleLineScreenStateAssertion {
	return NewSingleLineScreenStateAssertion(s, rowIndex, expectedOutput, fallbackPatterns, expectedPatternExplanation)
}

func (s *ScreenAsserter) WithPromptAssertion() *ScreenAsserter {
	s.PushAssertion(s.PromptAssertion("$ "))

	return s
}

func (s *ScreenAsserter) RunWithPromptAssertion() error {
	s.PushAssertion(s.PromptAssertion("$ "))
	defer s.PopAssertion()

	return s.Run()
}

func (s *ScreenAsserter) Run() error {
	currentRowIndex := 0
	// fmt.Println("Starting Run")
	// s.LogFullScreenState()

	for _, assertion := range s.Assertions {
		processedRowCount, err := assertion.Run(s.Shell.GetScreenState(), currentRowIndex)
		// fmt.Println(currentRowIndex, processedRowCount, assertion, err)
		if err != nil {
			return err
		}

		currentRowIndex += processedRowCount

		// Log "success" rows that were processed
		if s.nextRowToLog <= currentRowIndex {
			for i := s.nextRowToLog; i < currentRowIndex; i++ {
				s.LogRow(i)
			}
			s.nextRowToLog = currentRowIndex
		}
	}

	return nil
}

// ToDo: Remove in lieu of AsBool() helper
func (s *ScreenAsserter) RunBool() bool {
	// True if the prompt assertion is a success
	return s.Run() == nil
}

// Composition of Assertions

func (s *ScreenAsserter) PushAssertion(assertion Assertion) {
	s.Assertions = append(s.Assertions, assertion)
}

func (s *ScreenAsserter) PopAssertion() Assertion {
	if len(s.Assertions) == 0 {
		return nil
	}
	lastAssertion := s.Assertions[len(s.Assertions)-1]
	s.Assertions = s.Assertions[:len(s.Assertions)-1]
	return lastAssertion
}

func (s *ScreenAsserter) ClearAssertions() {
	s.Assertions = []Assertion{}
}
