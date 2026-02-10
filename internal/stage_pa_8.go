package internal

import (
	"fmt"
	"os"

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
	for i, word := range randomWords {
		entryName := prefix + word
		prefix = entryName + "_"

		// On the last run
		if i == len(randomWords)-1 {
			entryName += ".txt"
		}

		entryNames = append(entryNames, entryName)

		// Create a file at the end of the completion
		if i == len(randomWords)-1 {
			if err := writeFile(entryName, ""); err != nil {
				return err
			}
			defer func() {
				os.Remove(entryName)
			}()
			break
		}

		if err := os.Mkdir(entryName, 0755); err != nil {
			return err
		}
		defer func() {
			os.Remove(entryName)
		}()
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := random.RandomElementFromArray([]string{"ls", "stat", "file", "du"})

	initialTypedPrefix := fmt.Sprintf("%s %s", command, initialFilePrefix)
	reflections := []string{}

	for _, entryName := range entryNames {
		reflections = append(reflections, fmt.Sprintf("%s %s", command, entryName))
	}

	err := test_cases.CommandPartialCompletionsTestCase{
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
