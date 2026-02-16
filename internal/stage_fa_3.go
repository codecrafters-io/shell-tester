package internal

import (
	"fmt"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testFA3(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	// Create a directory to autocomplete to
	targetDirName := fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	targetDirPath := filepath.Join(workingDirPath, targetDirName)

	if err := MkdirAllWithTeardown(stageHarness, targetDirPath, 0755); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := "cd"

	typedPrefix := fmt.Sprintf("%s %s", command, targetDirName[:len(targetDirName)/2])
	// Directory completion should add a trailing slash
	completion := fmt.Sprintf("%s %s/", command, targetDirName)

	err = test_cases.AutocompleteTestCase{
		RawInput:           typedPrefix,
		ExpectedReflection: completion,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
