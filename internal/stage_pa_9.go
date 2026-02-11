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

func testPA9(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	words := random.RandomElementsFromArray(SMALL_WORDS, 6)
	nestedDirPath := filepath.Join("/tmp", words[0], words[1])
	if err := MkdirAllWithTeardown(stageHarness, nestedDirPath, 0755); err != nil {
		return err
	}

	nestedDirBaseName := filepath.Base(nestedDirPath)
	executable1BaseName := words[3]
	executable2BaseName := words[4]
	nonExecutableBaseName := words[5]

	topLevelDir := filepath.Join("/tmp", words[0])

	if err := WriteFilesWithTearDown(stageHarness, []WriteFileSpec{
		{
			FilePath:    filepath.Join(topLevelDir, executable1BaseName),
			FileContent: "",
			Permission:  0777,
		},
		{
			FilePath:    filepath.Join(topLevelDir, nonExecutableBaseName),
			FileContent: "",
			Permission:  0644,
		},
		{
			FilePath:    filepath.Join(nestedDirPath, executable2BaseName),
			FileContent: "",
			Permission:  0777,
		},
	}); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	allCompletions := []string{executable1BaseName, fmt.Sprintf("%s/", nestedDirBaseName)}
	slices.Sort(allCompletions)
	typedPrefix := fmt.Sprintf("%s/", topLevelDir)

	err := test_cases.MultipleCompletionsTestCase{
		RawPrefix:          typedPrefix,
		TabCount:           2,
		ExpectedReflection: strings.Join(allCompletions, "  "),
		ExpectedReflectionFallbackPatterns: []string{
			"^" + strings.Join(allCompletions, `\s*`) + "$",
		},
		SuccessMessage: fmt.Sprintf("Received completion for %q", typedPrefix),
		ExpectedAutocompletedReflectionHasNoSpace: true,
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, logger)

	if err != nil {
		return err
	}

	// We remove the assertion for the second line since it will be now used by newly typed prefix
	// instead of the old autocomplete
	// asserter.PopAssertion()

	// executable2Path := filepath.Join(nestedDirPath, executable2BaseName)
	// nestedDirIncompleteBaseName := nestedDirBaseName[:len(nestedDirBaseName)/2]

	// err = test_cases.PartialCompletionsTestCase{
	// 	ExistingPrefixInPromptLine:       typedPrefix,
	// 	Inputs:                           []string{nestedDirIncompleteBaseName, ""},
	// 	ExpectedReflections:              []string{fmt.Sprintf("%s/", nestedDirPath), executable2Path},
	// 	SuccessMessage:                   fmt.Sprintf("Received path completion for %q", typedPrefix),
	// 	SkipPromptAssertion:              true,
	// 	ExpectedLastReflectionHasNoSpace: false,
	// }.Run(asserter, shell, logger)

	// if err != nil {
	// 	return err
	// }

	return nil
}
