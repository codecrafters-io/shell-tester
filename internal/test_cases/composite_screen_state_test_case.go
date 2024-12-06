package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/tester-utils/logger"
)

// Stack-based implementation of CompositeScreenStateTestCase
// This is a stack of assertions that will be run in order
// It's a stack because we want to be able to push and pop assertions
type CompositeScreenStateTestCase struct {
	ScreenStateAssertions []assertions.Assertion
}

func NewCompositeScreenStateTestCase() *CompositeScreenStateTestCase {
	return &CompositeScreenStateTestCase{ScreenStateAssertions: []assertions.Assertion{}}
}

func (t *CompositeScreenStateTestCase) Run(screenState [][]string, logger *logger.Logger) error {
	// currentRowIndex := 0
	for _, assertion := range t.ScreenStateAssertions {
		if err := assertion.Run(screenState, logger); err != nil {
			return err
		}
	}

	return nil
}

func (t *CompositeScreenStateTestCase) PushAssertion(assertion assertions.Assertion) {
	t.ScreenStateAssertions = append(t.ScreenStateAssertions, assertion)
}

func (t *CompositeScreenStateTestCase) PopAssertion() assertions.Assertion {
	if len(t.ScreenStateAssertions) == 0 {
		panic("CodeCrafters Internal Error: no assertions to pop")
	}

	assertion := t.ScreenStateAssertions[len(t.ScreenStateAssertions)-1]
	t.ScreenStateAssertions = t.ScreenStateAssertions[:len(t.ScreenStateAssertions)-1]
	return assertion
}

func (t *CompositeScreenStateTestCase) PeekAssertion() assertions.Assertion {
	if len(t.ScreenStateAssertions) == 0 {
		panic("CodeCrafters Internal Error: no assertions to peek")
	}

	return t.ScreenStateAssertions[len(t.ScreenStateAssertions)-1]
}

func (t *CompositeScreenStateTestCase) AssertionsCount() int {
	return len(t.ScreenStateAssertions)
}
