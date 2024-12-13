package assertion_collection

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
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

func (c *AssertionCollection) RunWithPromptAssertion(screenState [][]string) *assertions.AssertionError {
	return c.runWithExtraAssertions(screenState, []assertions.Assertion{
		assertions.PromptAssertion{ExpectedPrompt: "$ "},
	})
}

func (c *AssertionCollection) RunWithoutPromptAssertion(screenState [][]string) *assertions.AssertionError {
	return c.runWithExtraAssertions(screenState, nil)
}

// ToDo: Remove all debug logs
func (c *AssertionCollection) runWithExtraAssertions(screenState [][]string, extraAssertions []assertions.Assertion) *assertions.AssertionError {
	allAssertions := append(c.Assertions, extraAssertions...)
	currentRowIndex := 0

	for _, assertion := range allAssertions {
		if len(screenState) == 0 {
			panic("CodeCrafters internal error: expected screen to have at least one row, but it was empty")
		}

		if currentRowIndex >= len(screenState) {
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
