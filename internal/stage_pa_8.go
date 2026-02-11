package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA8(stageHarness *test_case_harness.TestCaseHarness) error {
	stageLogger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	prefix := "xyz_"
	initialFilePrefix := prefix
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	entryNames := []string{}

	// Create the directories
	for _, word := range randomWords[:len(randomWords)-1] {
		dirName := prefix + word
		prefix = dirName + "_"

		if err := MkdirAllWithTeardown(stageHarness, dirName, 0755); err != nil {
			return err
		}

		entryNames = append(entryNames, dirName)
	}

	// Create a file for a completion
	fileName := fmt.Sprintf("%s%s.txt", prefix, randomWords[len(randomWords)-1])

	if err := WriteFileWithTeardown(stageHarness, fileName, "", 0644); err != nil {
		return err
	}

	// Add filename to entry names
	entryNames = append(entryNames, fileName)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := random.RandomElementFromArray([]string{"ls", "stat", "file", "du"})

	initialTypedPrefix := fmt.Sprintf("%s %s", command, initialFilePrefix)
	reflections := []string{}

	for _, entryName := range entryNames {
		reflections = append(reflections, fmt.Sprintf("%s %s", command, entryName))
	}

	err := test_cases.PartialCompletionsTestCase{
		Inputs:              []string{initialTypedPrefix, "_", "_"},
		ExpectedReflections: reflections,
		SuccessMessage:      fmt.Sprintf("Received all partial completions for %q", initialTypedPrefix),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageLogger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
