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

func testFA1(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	targetFileBaseName := fmt.Sprintf("%s-%d.txt", random.RandomWord(), random.RandomInt(1, 100))
	targetFilePath := filepath.Join(workingDirPath, targetFileBaseName)

	if err := WriteFileWithTeardown(stageHarness, targetFilePath, "", 0644); err != nil {
		return err
	}

	MustLogDirTree(stageHarness.Logger, workingDirPath)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	typedPrefix := fmt.Sprintf("%s %s", command, targetFileBaseName[:len(targetFileBaseName)/2])
	completion := fmt.Sprintf("%s %s ", command, targetFileBaseName)

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
