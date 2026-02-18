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

func testFA2(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	entryNames := random.RandomWords(3)
	targetFileRelativePath := fmt.Sprintf("%s/%s/%s.txt", entryNames[0], entryNames[1], entryNames[2])
	targetFileBaseName := filepath.Base(targetFileRelativePath)
	targetFileDirRelativePath := filepath.Dir(targetFileRelativePath)
	targetFileFullPath := filepath.Join(workingDirPath, targetFileRelativePath)

	// Create a nested file
	if err := WriteFileWithTeardown(stageHarness, targetFileFullPath, "", 0644); err != nil {
		return err
	}

	MustLogWorkingDirTree(stageHarness.Logger, workingDirPath)

	// Start and assert prompt
	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	typedPrefix := fmt.Sprintf("%s %s/", command, targetFileDirRelativePath)
	completion := fmt.Sprintf("%s %s ", command, filepath.Join(targetFileDirRelativePath, targetFileBaseName))

	err = test_cases.AutocompleteTestCase{
		RawInput:            typedPrefix,
		ExpectedCompletion:  completion,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageHarness.Logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
