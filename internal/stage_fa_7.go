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

func testFA7(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	shell.SetWorkingDirectory(workingDirPath)
	prefix := fmt.Sprintf("%s_", random.RandomElementFromArray(SMALL_WORDS))
	suffixes := random.RandomInts(1, 10, 2)
	fileSuffix := suffixes[0]
	dirSuffix := suffixes[1]

	fileBaseName := fmt.Sprintf("%s%d", prefix, fileSuffix)
	dirBaseName := fmt.Sprintf("%s%d", prefix, dirSuffix)
	if err := WriteFileWithTeardown(stageHarness, filepath.Join(workingDirPath, fileBaseName), "", 0644); err != nil {
		return err
	}
	if err := MkdirAllWithTeardown(stageHarness, filepath.Join(workingDirPath, dirBaseName), 0755); err != nil {
		return err
	}

	allCompletions := []string{fileBaseName, dirBaseName + "/"}
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

	// Test file completion for first argument
	expectedReflectionAfterFileCompletion := fmt.Sprintf("%s %s", command, fileBaseName)

	err = test_cases.AutocompleteTestCase{
		PreExistingInputOnLine: initialTypedPrefix,
		RawInput:               fmt.Sprintf("%d", fileSuffix),
		ExpectedReflection:     expectedReflectionAfterFileCompletion,
		ExpectedAutocompletedReflectionHasNoSpace: false,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	// Complete to directory for the second argument
	expectedReflectionAfterDirCompletion := fmt.Sprintf("%s %s %s/", command, fileBaseName, dirBaseName)
	err = test_cases.AutocompleteTestCase{
		// The extra space should be inserted by previous step
		PreExistingInputOnLine: fmt.Sprintf("%s ", expectedReflectionAfterFileCompletion),
		RawInput:               dirBaseName,
		ExpectedReflection:     expectedReflectionAfterDirCompletion,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	// Check invalid completion now
	invalidCompletionRawInput := fmt.Sprintf(" missing_entry-%d", random.RandomInt(1, 1000))
	expectedReflectionAfterInvalidCompletion := fmt.Sprintf("%s%s", expectedReflectionAfterDirCompletion, invalidCompletionRawInput)
	err = test_cases.AutocompleteTestCase{
		PreExistingInputOnLine: expectedReflectionAfterDirCompletion,
		RawInput:               invalidCompletionRawInput,
		ExpectedReflection:     expectedReflectionAfterInvalidCompletion,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
		CheckForBell:                              true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)

	// TODO: Kinda unsatisfied by the logs produced for this stage:
	// We can't log intermediate screen state in the middle: is confusing/misleading
	// Should we perform shell teardown and test separately for three different arguments 1st, 2nd and 3rd separately
	// But keep the scenario same as before?
}
