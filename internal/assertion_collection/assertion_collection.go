package assertion_collection

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

const ShouldPrintDebugLogs = false

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

	if ShouldPrintDebugLogs {
		printScreenState(screenState)
	}

	for _, assertion := range allAssertions {
		processedRowCount, err := assertion.Run(screenState, currentRowIndex)
		if err != nil {
			if ShouldPrintDebugLogs {
				fmt.Printf("❌ %s\n", assertion.Inspect())
			}

			return err
		}

		if ShouldPrintDebugLogs {
			fmt.Printf("✅ %s (%d rows) currentRowIndex: %d\n", assertion.Inspect(), processedRowCount, currentRowIndex)
		}

		if c.OnAssertionSuccess != nil {
			c.OnAssertionSuccess(currentRowIndex, processedRowCount)
		}

		currentRowIndex += processedRowCount
	}

	return nil
}

func printScreenState(screenState [][]string) {
	fmt.Println("--- Screen start ---")

	for _, row := range screenState {
		cleanedRow := utils.BuildCleanedRow(row)

		if len(cleanedRow) != 0 {
			fmt.Println(cleanedRow)
		}
	}

	fmt.Println("--- Screen end ----")
}
