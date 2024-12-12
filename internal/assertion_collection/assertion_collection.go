package assertion_collection

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
)

type AssertionCollection struct {
	Assertions []assertions.Assertion
}

func NewAssertionCollection() *AssertionCollection {
	return &AssertionCollection{Assertions: []assertions.Assertion{}}
}

func (c *AssertionCollection) AddAssertion(assertion assertions.Assertion) {
	c.Assertions = append(c.Assertions, assertion)
}

func (c *AssertionCollection) RunWithPromptAssertion(screenState [][]string) error {
	return c.runWithExtraAssertions(screenState, []assertions.Assertion{
		assertions.PromptAssertion{ExpectedPrompt: "$ "},
	})
}

func (c *AssertionCollection) RunWithoutPromptAssertion(screenState [][]string) error {
	return c.runWithExtraAssertions(screenState, nil)
}

func (c *AssertionCollection) runWithExtraAssertions(screenState [][]string, extraAssertions []assertions.Assertion) error {
	allAssertions := append(c.Assertions, extraAssertions...)
	currentRowIndex := 0

	for _, assertion := range allAssertions {
		processedRowCount, err := assertion.Run(screenState, currentRowIndex)
		if err != nil {
			return err
		}

		currentRowIndex += processedRowCount
	}

	return nil
}
