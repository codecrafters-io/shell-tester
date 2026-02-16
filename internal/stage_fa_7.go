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

func testFA7(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	// Create two files for multi-argument completion test
	file1BaseName := fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	file2BaseName := fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	file1Path := filepath.Join(workingDirPath, file1BaseName)
	file2Path := filepath.Join(workingDirPath, file2BaseName)

	if err := WriteFileWithTeardown(stageHarness, file1Path, "", 0644); err != nil {
		return err
	}
	if err := WriteFileWithTeardown(stageHarness, file2Path, "", 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	// First argument: complete file1
	typedPrefix1 := fmt.Sprintf("%s %s", command, file1BaseName[:len(file1BaseName)/2])
	completion1 := fmt.Sprintf("%s %s", command, file1BaseName)

	err = test_cases.AutocompleteTestCase{
		RawInput:           typedPrefix1,
		ExpectedReflection: completion1,
		ExpectedAutocompletedReflectionHasNoSpace: false,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	// Second argument: complete file2
	// Type partial file2 name after the space added by first completion
	typedPrefix2 := file2BaseName[:len(file2BaseName)/2]
	completion2 := fmt.Sprintf("%s %s %s", command, file1BaseName, file2BaseName)

	err = test_cases.AutocompleteTestCase{
		RawInput:           typedPrefix2,
		ExpectedReflection: completion2,
		ExpectedAutocompletedReflectionHasNoSpace: false,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
