package internal

import (
	"fmt"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA7(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	prefix := "pear_"
	suffixes := random.RandomInts(1, 10, 3)

	fileSuffixes := suffixes[:2]
	dirSuffixes := suffixes[2:]

	allCompletions := []string{}

	for _, integerSuffix := range fileSuffixes {
		fileName := fmt.Sprintf("%s%d", prefix, integerSuffix)

		if err := WriteFileWithTeardown(stageHarness, fileName, "", 0644); err != nil {
			return err
		}

		allCompletions = append(allCompletions, fileName)
	}

	for _, integerSuffix := range dirSuffixes {
		dirName := fmt.Sprintf("%s%d", prefix, integerSuffix)

		if err := MkdirAllWithTeardown(stageHarness, dirName, 0755); err != nil {
			return err
		}

		allCompletions = append(allCompletions, fmt.Sprintf("%s/", dirName))
	}

	slices.Sort(allCompletions)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	commands := []string{"ls", "cat", "stat", "file"}
	typedPrefix := fmt.Sprintf("%s %s", random.RandomElementFromArray(commands), prefix)

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

	return logAndQuit(asserter, nil)
}
