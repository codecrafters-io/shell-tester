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

func testFA2(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	// Create a nested directory with a file
	nestedDirName := random.RandomWord()
	targetFileBaseName := fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	nestedDirPath := filepath.Join(workingDirPath, nestedDirName)
	targetFilePath := filepath.Join(nestedDirPath, targetFileBaseName)

	if err := WriteFileWithTeardown(stageHarness, targetFilePath, "", 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	// Type partial path: "nestedDir/partialFile"
	relativeFilePath := filepath.Join(nestedDirName, targetFileBaseName)
	typedPrefix := fmt.Sprintf("%s %s", command, relativeFilePath[:len(nestedDirName)+1+len(targetFileBaseName)/2])
	completion := fmt.Sprintf("%s %s", command, relativeFilePath)

	err = test_cases.AutocompleteTestCase{
		RawInput:           typedPrefix,
		ExpectedReflection: completion,
		ExpectedAutocompletedReflectionHasNoSpace: false,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
