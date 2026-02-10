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

func testPA3(stageHarness *test_case_harness.TestCaseHarness) error {
	filePath, _, err := GetRandomFile(stageHarness, "txt", 0644)
	if err != nil {
		return err
	}

	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	randomCommand := random.RandomElementFromArray([]string{
		"cat",
		"ls",
		"stat",
	})

	fileDir := filepath.Dir(filePath)
	fileBaseName := filepath.Base(filePath)
	filePartialPath := filepath.Join(fileDir, fileBaseName[:len(fileBaseName)/2])

	typedPrefix := fmt.Sprintf("%s %s", randomCommand, filePartialPath)
	completion := fmt.Sprintf("%s %s", randomCommand, filePath)

	err = test_cases.AutocompleteTestCase{
		TypedPrefix:                               typedPrefix,
		ExpectedPromptLineReflection:              completion,
		ExpectedAutocompletedReflectionHasNoSpace: false,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
