package assertion_collection

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

type AssertionCollection struct {
	Assertions []assertions.Assertion

	OnAssertionSuccess func(startRowIndex int, processedRowCount int)
}

func NewAssertionCollection() *AssertionCollection {
	return &AssertionCollection{Assertions: []assertions.Assertion{}}
}

func (c *AssertionCollection) AddAssertion(assertion assertions.Assertion) {
	c.Assertions = append(c.Assertions, assertion)
}

func (c *AssertionCollection) PopAssertion() assertions.Assertion {
	if len(c.Assertions) == 0 {
		panic("CodeCrafters internal error: no assertions to pop")
	}

	lastAssertion := c.Assertions[len(c.Assertions)-1]
	c.Assertions = c.Assertions[:len(c.Assertions)-1]
	return lastAssertion
}

func (c *AssertionCollection) RunWithPromptAssertion(screenState screen_state.ScreenState) *assertions.AssertionError {
	return c.runWithExtraAssertions(screenState, []assertions.Assertion{
		assertions.PromptAssertion{ExpectedPrompt: utils.PROMPT},
	})
}

func (c *AssertionCollection) RunWithoutPromptAssertion(screenState screen_state.ScreenState) *assertions.AssertionError {
	return c.runWithExtraAssertions(screenState, nil)
}

func (c *AssertionCollection) runWithExtraAssertions(screenState screen_state.ScreenState, extraAssertions []assertions.Assertion) *assertions.AssertionError {
	allAssertions := append(c.Assertions, extraAssertions...)
	currentRowIndex := 0

	for _, assertion := range allAssertions {
		if screenState.GetRowCount() == 0 {
			panic("CodeCrafters internal error: expected screen to have at least one row, but it was empty")
		}

		if currentRowIndex > screenState.GetRowCount()-1 {
			panic("CodeCrafters internal error: startRowIndex is larger than screenState rows")
		}

		processedRowCount, err := assertion.Run(screenState, currentRowIndex)
		if err != nil {
			return err
		}

		if c.OnAssertionSuccess != nil {
			c.OnAssertionSuccess(currentRowIndex, processedRowCount)
		}

		currentRowIndex += processedRowCount
	}

	return nil
}
