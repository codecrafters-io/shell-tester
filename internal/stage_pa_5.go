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

func testPA5(stageHarness *test_case_harness.TestCaseHarness) error {
	stageLogger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	fileDirPath, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	fileBaseName, _, err := CreateRandomFileInDir(stageHarness, fileDirPath, "txt", 0644)
	if err != nil {
		return err
	}

	filePath := filepath.Join(fileDirPath, fileBaseName)
	fileDirParentPath := filepath.Dir(fileDirPath)
	fileDirParentBaseName := filepath.Base(fileDirParentPath)
	fileDirGrandParentPath := filepath.Dir(fileDirParentPath)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := random.RandomElementFromArray([]string{"ls", "stat", "file", "du"})

	initialTypedPrefix := fmt.Sprintf(
		"%s %s",
		command,
		filepath.Join(fileDirGrandParentPath, fileDirParentBaseName[:len(fileDirParentBaseName)/2]),
	)

	reflections := []string{
		fmt.Sprintf("%s %s/", command, fileDirParentPath),
		fmt.Sprintf("%s %s/", command, fileDirPath),
		fmt.Sprintf("%s %s", command, filePath),
	}

	err = test_cases.PartialCompletionsTestCase{
		Inputs:              []string{initialTypedPrefix, "", ""},
		ExpectedReflections: reflections,
		SuccessMessage:      fmt.Sprintf("Received all path completions for %q", initialTypedPrefix),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageLogger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
