package test_cases

import "github.com/codecrafters-io/shell-tester/internal/assertions"

type CompositeScreenStateTestCase struct {
	ScreenStateAssertions []assertions.Assertion
}
