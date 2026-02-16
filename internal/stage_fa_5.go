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

func testFA5(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	shell.SetWorkingDirectory(workingDirPath)
	prefix := fmt.Sprintf("%s_", random.RandomElementFromArray(SMALL_WORDS))
	fileSuffix := random.RandomInt(1, 10)
	dirSuffix := random.RandomInt(1, 10)
	for dirSuffix == fileSuffix {
		dirSuffix = random.RandomInt(1, 10)
	}

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
		RawInput:           initialTypedPrefix,
		TabCount:           2,
		ExpectedCompletion: strings.Join(allCompletions, "  "),
		ExpectedCompletionOptionsFallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile("^" + strings.Join(allCompletions, `\s*`) + "$"),
		},
		SuccessMessage:      fmt.Sprintf("Received completion for %q", initialTypedPrefix),
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)
	if err != nil {
		return err
	}

	// Complete to either the file or the dir
	var completionSuffix int
	var expectedCompletion string
	if random.RandomInt(0, 2) == 0 {
		completionSuffix = fileSuffix
		expectedCompletion = fmt.Sprintf("%s %s%d", command, prefix, fileSuffix)
	} else {
		completionSuffix = dirSuffix
		expectedCompletion = fmt.Sprintf("%s %s%d/", command, prefix, dirSuffix)
	}

	err = test_cases.AutocompleteTestCase{
		PreviousInputOnLine:          initialTypedPrefix,
		RawInput:                     fmt.Sprintf("%d", completionSuffix),
		ExpectedCompletion:           expectedCompletion,
		ExpectedCompletionHasNoSpace: (completionSuffix == dirSuffix),
		SkipPromptAssertion:          true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
