package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testFA4(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	// Use a non-existent file name to trigger the bell
	nonExistentFile := fmt.Sprintf("nonexistent_%s_%d", random.RandomWord(), random.RandomInt(1, 100))
	typedPrefix := fmt.Sprintf("%s %s", command, nonExistentFile)

	// When there's no match, the shell should ring a bell and keep the input unchanged
	err = test_cases.AutocompleteTestCase{
		RawInput:           typedPrefix,
		ExpectedReflection: typedPrefix,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
