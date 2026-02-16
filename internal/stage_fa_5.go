package internal

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testFA5(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	// Create multiple files with a common prefix
	prefix := fmt.Sprintf("file_%d_", random.RandomInt(1, 100))
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	fileNames := []string{}
	for _, word := range randomWords {
		fileName := prefix + word
		fileNames = append(fileNames, fileName)
		filePath := filepath.Join(workingDirPath, fileName)
		if err := WriteFileWithTeardown(stageHarness, filePath, "", 0644); err != nil {
			return err
		}
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()
	typedPrefix := fmt.Sprintf("%s %s", command, prefix)

	sort.Strings(fileNames)
	completions := strings.Join(fileNames, "  ")

	err = test_cases.MultipleCompletionsTestCase{
		RawInput:           typedPrefix,
		TabCount:           2,
		ExpectedReflection: completions,
		SuccessMessage:     fmt.Sprintf("Received completions for %q", prefix),
		ExpectedAutocompletedReflectionHasNoSpace: true,
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
