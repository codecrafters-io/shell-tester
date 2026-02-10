package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA6(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	randomCommand := random.RandomElementFromArray([]string{
		"cat",
		"less",
		"ls",
		"stat",
	})

	typedPrefixInteger := random.RandomInt(1, 1000)
	typedPrefix := fmt.Sprintf("%s missing_%d", randomCommand, typedPrefixInteger)

	err := test_cases.AutocompleteTestCase{
		TypedPrefix:                               typedPrefix,
		ExpectedPromptLineReflection:              typedPrefix,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
		CheckForBell:                              true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
