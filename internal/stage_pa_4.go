package internal

import (
	"fmt"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA4(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	dirPath, err := GetRandomDirectory(stageHarness)

	if err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	randomCommand := random.RandomElementFromArray([]string{
		"cat",
		"ls",
		"stat",
	})

	dirParentPath := filepath.Dir(dirPath)
	dirBaseName := filepath.Base(dirPath)
	dirPartialPath := filepath.Join(dirParentPath, dirBaseName[:len(dirBaseName)/2])

	typedPrefix := fmt.Sprintf("%s %s", randomCommand, dirPartialPath)
	completion := fmt.Sprintf("%s %s/", randomCommand, dirPath)

	err = test_cases.AutocompleteTestCase{
		TypedPrefix:                               typedPrefix,
		ExpectedPromptLineReflection:              completion,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
