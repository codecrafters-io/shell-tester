package internal

import (
	"fmt"
	"path/filepath"
	"slices"
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
	prefix := fmt.Sprintf("%s_", random.RandomElementFromArray(SMALL_WORDS))
	suffixes := random.RandomInts(1, 10, 3)

	// Create two files and one dir
	fileSuffixes := suffixes[:2]
	dirSuffixes := suffixes[2:]

	allCompletions := []string{}

	// Create files
	for _, fileSuffix := range fileSuffixes {
		fileBaseName := fmt.Sprintf("%s%d", prefix, fileSuffix)
		filePath := filepath.Join(workingDirPath, fileBaseName)
		allCompletions = append(allCompletions, fileBaseName)

		if err := WriteFileWithTeardown(stageHarness, filePath, "", 0644); err != nil {
			return err
		}
	}

	// Create directories
	for _, dirSuffix := range dirSuffixes {
		dirbaseName := fmt.Sprintf("%s%d", prefix, dirSuffix)
		dirPath := filepath.Join(workingDirPath, dirbaseName)
		allCompletions = append(allCompletions, fmt.Sprintf("%s/", dirbaseName))

		if err := MkdirAllWithTeardown(stageHarness, dirPath, 0755); err != nil {
			return err
		}
	}

	slices.Sort(allCompletions)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFileAndDir()
	initialTypedPrefix := fmt.Sprintf("%s %s", command, prefix)

	err = test_cases.MultipleCompletionsTestCase{
		RawInput:                  initialTypedPrefix,
		TabCount:                  2,
		ExpectedCompletionOptions: strings.Join(allCompletions, "  "),
		ExpectedCompletionOptionsFallbackPatterns: []string{
			"^" + strings.Join(allCompletions, `\s*`) + "$",
		},
		SuccessMessage:      fmt.Sprintf("Received completion for %q", initialTypedPrefix),
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)
	if err != nil {
		return err
	}

	// For next test case, we check for new input reflection instead of the old one
	asserter.PopAssertion()

	suffix := random.RandomElementFromArray(suffixes)

	expectedReflection := fmt.Sprintf("%s %s%d", command, prefix, suffix)

	// Reflection should contain space if the matched entry was a file, trailing slash for dir
	if slices.Contains(dirSuffixes, suffix) {
		expectedReflection += "/"
	} else {
		expectedReflection += " "
	}

	err = test_cases.AutocompleteTestCase{
		PreExistingInputOnLine: initialTypedPrefix,
		RawInput:               fmt.Sprintf("%d", suffix),
		ExpectedReflection:     expectedReflection,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
