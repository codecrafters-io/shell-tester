package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	executableDir, err := SetUpCustomCommands(stageHarness, shell, []string{"cat"})
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	defer cleanupDirectories([]string{executableDir})

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
	}

	executableName1 := `'exe  with  space'`
	executableName2 := `'exe with "quotes"'`
	executableName3 := `"exe with \'single quotes\'"`
	executableName4 := `'exe with \n newline'`

	originalExecutablePath := "/tmp/custom_cat_executable"
	err = createExecutableCallingCat(originalExecutablePath)
	if err != nil {
		panic("CodeCrafters Internal Error: Cannot create executable")
	}

	err = custom_executable.CopyFileToMultiplePaths(originalExecutablePath, []string{path.Join(executableDir, executableName1), path.Join(executableDir, executableName2), path.Join(executableDir, executableName3), path.Join(executableDir, executableName4)}, logger)
	if err != nil {
		panic("CodeCrafters Internal Error: Cannot copy executable")
	}

	writeFiles([]string{
		path.Join(executableDir, "f1"),
		path.Join(executableDir, "f2"),
		path.Join(executableDir, "f3"),
		path.Join(executableDir, "f4"),
	}, []string{fileContents[0] + "\n", fileContents[1] + "\n", fileContents[2] + "\n", fileContents[3] + "\n"}, logger)

	inputs := []string{
		fmt.Sprintf(`%s %s/f1`, executableName1, executableDir),
		fmt.Sprintf(`%s %s/f2`, executableName2, executableDir),
		fmt.Sprintf(`%s %s/f3`, executableName3, executableDir),
		fmt.Sprintf(`%s %s/f4`, executableName4, executableDir),
	}
	expectedOutputs := []string{
		fileContents[0],
		fileContents[1],
		fileContents[2],
		fileContents[3],
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.CommandResponseTestCase{
			Command:          testCaseContent.Input,
			ExpectedOutput:   testCaseContent.ExpectedOutput,
			FallbackPatterns: nil,
			SuccessMessage:   "✓ Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	return logAndQuit(asserter, nil)
}

func createExecutableFile(path string, contents string) error {
	return os.WriteFile(path, []byte(contents), 0o755)
}

func createExecutableCallingCat(path string) error {
	content := `#!/bin/sh
exec cat "$@"`

	return createExecutableFile(path, content)
}
