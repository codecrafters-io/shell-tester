package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, currentPath))

	if err := shell.Start(); err != nil {
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

	originalExecutablePath := "/tmp/catexe"
	err = createExecutableCallingCat(originalExecutablePath)
	if err != nil {
		panic("CodeCrafters Internal Error: Cannot create executable")
	}

	err = custom_executable.CopyExecutableToMultiplePaths(originalExecutablePath, []string{path.Join(randomDir, executableName1), path.Join(randomDir, executableName2), path.Join(randomDir, executableName3), path.Join(randomDir, executableName4)}, logger)
	if err != nil {
		panic("CodeCrafters Internal Error: Cannot copy executable")
	}

	writeFiles([]string{
		path.Join(randomDir, "f1"),
		path.Join(randomDir, "f2"),
		path.Join(randomDir, "f3"),
		path.Join(randomDir, "f4"),
	}, []string{fileContents[0] + "\n", fileContents[1] + "\n", fileContents[2] + "\n", fileContents[3] + "\n"}, logger)

	inputs := []string{
		fmt.Sprintf(`%s %s/f1`, executableName1, randomDir),
		fmt.Sprintf(`%s %s/f2`, executableName2, randomDir),
		fmt.Sprintf(`%s %s/f3`, executableName3, randomDir),
		fmt.Sprintf(`%s %s/f4`, executableName4, randomDir),
	}
	expectedOutputs := []string{
		fileContents[0],
		fileContents[1],
		fileContents[2],
		fileContents[3],
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}

func createExecutableFile(path string, contents string) error {
	return os.WriteFile(path, []byte(contents), 0o755)
}

func createExecutableCallingCat(path string) error {
	content := `#!/bin/sh
exec cat "$@"`
	return createExecutableFile(path, content)
}
