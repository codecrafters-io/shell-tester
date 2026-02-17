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

	dirNames := random.RandomElementsFromArray(SMALL_WORDS, 2)
	dirRelativePath := filepath.Join(dirNames[0], dirNames[1])
	dirParentRelativePath := filepath.Dir(dirRelativePath)
	dirAbsolutePath := filepath.Join(workingDirPath, dirRelativePath)

	if err := MkdirAllWithTeardown(stageHarness, dirAbsolutePath, 0755); err != nil {
		return err
	}

	MustLogDirTree(stageHarness.Logger, workingDirPath)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForDir()
	initialPrefix := fmt.Sprintf("%s ", command)
	firstCompletion := fmt.Sprintf("%s %s/", command, dirParentRelativePath)
	lastCompletion := fmt.Sprintf("%s %s/", command, dirRelativePath)

	err = test_cases.PartialCompletionsTestCase{
		InputAndCompletionPairs: []test_cases.InputAndCompletionPair{
			{
				Input:              initialPrefix,
				ExpectedCompletion: firstCompletion,
			},
			{
				Input:              "",
				ExpectedCompletion: lastCompletion,
			},
		},
		SuccessMessage:      fmt.Sprintf("Received all path completions for %q", initialPrefix),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
