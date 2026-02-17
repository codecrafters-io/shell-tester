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

func testFA6(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	shell.SetWorkingDirectory(workingDirPath)

	prefix := "xyz_"
	initialFilePrefix := prefix
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	entryNames := []string{}

	// Create the directories
	for _, word := range randomWords[:len(randomWords)-1] {
		dirName := prefix + word
		prefix = dirName + "_"

		if err := MkdirAllWithTeardown(stageHarness, filepath.Join(workingDirPath, dirName), 0755); err != nil {
			return err
		}

		entryNames = append(entryNames, dirName)
	}

	// Create a file for a completion
	fileName := fmt.Sprintf("%s%s.txt", prefix, randomWords[len(randomWords)-1])

	if err := WriteFileWithTeardown(stageHarness, filepath.Join(workingDirPath, fileName), "", 0644); err != nil {
		return err
	}

	MustLogDirTree(stageHarness.Logger, workingDirPath)

	// Add filename to entry names
	entryNames = append(entryNames, fileName)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFileAndDir()
	initialTypedPrefix := fmt.Sprintf("%s %s", command, initialFilePrefix)
	expectedCompletions := []string{}

	for _, entryName := range entryNames {
		expectedCompletions = append(expectedCompletions, fmt.Sprintf("%s %s", command, entryName))
	}

	inputAndCompletionPairs := []test_cases.InputAndCompletionPair{}
	for idx, expectedCompletion := range expectedCompletions {
		input := "_"

		if idx == 0 {
			input = initialTypedPrefix
		}

		// Extra space at the end of final completion
		if idx == len(expectedCompletions)-1 {
			expectedCompletion += " "
		}

		inputAndCompletionPairs = append(inputAndCompletionPairs, test_cases.InputAndCompletionPair{
			Input:              input,
			ExpectedCompletion: expectedCompletion,
		})
	}

	err = test_cases.PartialCompletionsTestCase{
		InputAndCompletionPairs: inputAndCompletionPairs,
		SuccessMessage:          fmt.Sprintf("Received all partial completions for %q", initialTypedPrefix),
		SkipPromptAssertion:     true,
	}.Run(asserter, shell, stageHarness.Logger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
