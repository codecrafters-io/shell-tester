package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testA5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	prefix := "xyz_"
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	executableNames := []string{}
	for _, word := range randomWords {
		executableName := prefix + word
		executableNames = append(executableNames, executableName)
		_, err := SetUpSignaturePrinter(stageHarness, shell, getRandomString(), executableName)
		if err != nil {
			return err
		}
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := prefix
	sort.Strings(executableNames)
	completions := strings.Join(executableNames, "  ")
	completionEndsWithNoSpace := true

	err := test_cases.CommandMultipleCompletionsTestCase{
		RawCommand:         command,
		TabCount:           2,
		ExpectedReflection: completions,
		SuccessMessage:     fmt.Sprintf("Received completion for %q", command),
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

func logScreenState(shell *shell_executable.ShellExecutable) {
	screenState := shell.GetScreenState()
	for _, row := range screenState {
		cleanedRow := virtual_terminal.BuildCleanedRow(row)
		if len(cleanedRow) > 0 {
			fmt.Println(cleanedRow)
		}
	}
}
