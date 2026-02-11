package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA1(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	workingDirPath, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	shell.SetWorkingDirectory(workingDirPath)
	fileBaseName, _, err := CreateRandomFileInDir(stageHarness, workingDirPath, "txt", 0644)

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
		"du",
	})

	typedPrefix := fmt.Sprintf("%s %s", randomCommand, fileBaseName[:len(fileBaseName)/2])
	completion := fmt.Sprintf("%s %s", randomCommand, fileBaseName)

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
