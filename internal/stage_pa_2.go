package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA2(stageHarness *test_case_harness.TestCaseHarness) error {
	workingDirPath, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	directoryBaseName, err := CreateRandomSubDir(stageHarness, workingDirPath)
	if err != nil {
		return err
	}

	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.SetWorkingDirectory(workingDirPath)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := random.RandomElementFromArray([]string{
		"ls",
		"stat",
		"rmdir",
	})

	incompleteDirectoryName := directoryBaseName[:len(directoryBaseName)/2]
	typedPrefix := fmt.Sprintf("%s %s", command, incompleteDirectoryName)
	completion := fmt.Sprintf("%s %s/", command, directoryBaseName)

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
