package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	fileDir := "/tmp/"
	fileDir = filepath.Join(fileDir, random.RandomElementFromArray([]string{"foo", "bar", "baz"}))
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, 0755)
	}

	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", fileDir, currentPath))

	if err := shell.Start(); err != nil {
		return err
	}

	_, L := getRandomWordsSmallAndLarge(5, 6)
	adjectives := random.RandomElementsFromArray(ADJECTIVES, 4)

	file1Contents := fmt.Sprintf(`'%s'"'%s'`, L[0], L[1])
	file2Contents := fmt.Sprintf(`'%s'\"\"'%s'`, L[2], L[3])
	file3Contents := fmt.Sprintf(`'%s'\''%s'`, L[4], L[5])
	file4Contents := fmt.Sprintf(`%s\ncat`, adjectives[3])

	executableName1 := fmt.Sprintf(`'%s     cat'`, adjectives[0])
	executableName2 := fmt.Sprintf(`"\"%s\" cat"`, adjectives[1])
	executableName3 := fmt.Sprintf(`"'%s'"\ \ 'cat'`, adjectives[2])
	executableName4 := fmt.Sprintf(`%s\ncat`, adjectives[3])

	err := custom_executable.CopyExecutableToMultiplePaths("/usr/bin/cat", []string{path.Join(fileDir, executableName1), path.Join(fileDir, executableName2), path.Join(fileDir, executableName3), path.Join(fileDir, executableName4)}, logger)
	if err != nil {
		panic("CodeCrafters Internal Error: Cannot copy executable")
	}

	writeFiles([]string{
		path.Join(fileDir, "f1"),
		path.Join(fileDir, "f2"),
		path.Join(fileDir, "f3"),
		path.Join(fileDir, "f4"),
	}, []string{file1Contents + "\n", file2Contents + "\n", file3Contents + "\n", file4Contents + "\n"}, logger)

	inputs := []string{
		fmt.Sprintf(`%s %s/f1`, executableName1, fileDir),
		fmt.Sprintf(`%s %s/f2`, executableName2, fileDir),
		fmt.Sprintf(`%s %s/f3`, executableName3, fileDir),
		fmt.Sprintf(`%s %s/f4`, executableName4, fileDir),
	}
	expectedOutputs := []string{
		fmt.Sprintf(`'%s'"'%s'`, L[0], L[1]),
		fmt.Sprintf(`'%s'\"\"'%s'`, L[2], L[3]),
		fmt.Sprintf(`'%s'\''%s'`, L[4], L[5]),
		fmt.Sprintf(`%s\ncat`, adjectives[3]),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineStringMatchTestCase{
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
