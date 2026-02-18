package internal

import (
	"fmt"
	"path/filepath"
	"regexp"
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

	fileBaseName := fmt.Sprintf("%s%d.txt", prefix, fileSuffix)
	dirBaseName := fmt.Sprintf("%s%d", prefix, dirSuffix)
	if err := WriteFileWithTeardown(stageHarness, filepath.Join(workingDirPath, fileBaseName), "", 0644); err != nil {
		return err
	}
	if err := MkdirAllWithTeardown(stageHarness, filepath.Join(workingDirPath, dirBaseName), 0755); err != nil {
		return err
	}

	MustLogWorkingDirTree(stageHarness.Logger, workingDirPath)

	allCompletions := []string{fileBaseName, dirBaseName + "/"}
	slices.Sort(allCompletions)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFileAndDir()
	initialTypedPrefix := fmt.Sprintf("%s %s", command, prefix)

	escapedCompletions := make([]string, len(allCompletions))
	for i, c := range allCompletions {
		escapedCompletions[i] = regexp.QuoteMeta(c)
	}

	// Test multiple completions for first argument
	err = test_cases.MultipleCompletionsTestCase{
		RawInput:                      initialTypedPrefix,
		TabCount:                      2,
		ExpectedCompletionOptionsLine: strings.Join(allCompletions, "  "),
		ExpectedCompletionOptionsLineFallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile("^" + strings.Join(escapedCompletions, `\s+`) + "$"),
		},
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)
	if err != nil {
		return err
	}

	// Test file completion for first argument
	expectedCompletionAfterFileCompletion := fmt.Sprintf("%s %s ", command, fileBaseName)

	// Type some prefix so that the tab press will result in autocompletion of the first argument
	err = test_cases.AutocompleteTestCase{
		PreviousInputOnLine: initialTypedPrefix,
		RawInput:            fmt.Sprintf("%d", fileSuffix),
		ExpectedCompletion:  expectedCompletionAfterFileCompletion,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	// Complete to directory for the second argument
	expectedCompletionAfterDirCompletion := fmt.Sprintf("%s %s %s/", command, fileBaseName, dirBaseName)
	err = test_cases.AutocompleteTestCase{
		PreviousInputOnLine: expectedCompletionAfterFileCompletion,
		RawInput:            dirBaseName,
		ExpectedCompletion:  expectedCompletionAfterDirCompletion,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	// Check invalid completion now
	// Insert extra space at the beginning to separate arguments since dir completion does not put an extra
	// space at the end of the completion
	invalidCompletionRawInput := fmt.Sprintf(" missing_entry-%d", random.RandomInt(1, 1000))
	expectedFinalCompletion := fmt.Sprintf("%s%s", expectedCompletionAfterDirCompletion, invalidCompletionRawInput)
	err = test_cases.AutocompleteTestCase{
		PreviousInputOnLine: expectedCompletionAfterDirCompletion,
		RawInput:            invalidCompletionRawInput,
		ExpectedCompletion:  expectedFinalCompletion,
		SkipPromptAssertion: true,
		CheckForBell:        true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
