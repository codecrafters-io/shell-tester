package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	fileDir := "/tmp/"
	fileDir = filepath.Join(fileDir, random.RandomElementFromArray([]string{"foo", "bar", "baz"}))
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, 0755)
	}

	_, L := getRandomWordsSmallAndLarge(5, 6)
	file1Contents := fmt.Sprintf(`'%s'"'%s'`, L[0], L[1])
	file2Contents := fmt.Sprintf(`'%s'\"\"'%s'`, L[2], L[3])
	file3Contents := fmt.Sprintf(`'%s'\''%s'`, L[4], L[5])

	writeFiles([]string{
		path.Join(fileDir, "f1"),
		path.Join(fileDir, "f2"),
		path.Join(fileDir, "f3"),
	}, []string{file1Contents, file2Contents, file3Contents + "\n"}, logger)

	inputs := []string{
		fmt.Sprintf(`echo '%s\\\n%s'`, L[0], L[1]),
		fmt.Sprintf(`echo '%s\"%s%s\"%s'`, L[0], L[1], L[2], L[3]),
		fmt.Sprintf(`echo '%s\\n%s'`, L[2], L[3]),
		fmt.Sprintf(`cat %s/f1 %s/f2 %s/f3`, fileDir, fileDir, fileDir),
	}
	expectedOutputs := []string{
		fmt.Sprintf(`%s\\\n%s`, L[0], L[1]),
		fmt.Sprintf(`%s\"%s%s\"%s`, L[0], L[1], L[2], L[3]),
		fmt.Sprintf(`%s\\n%s`, L[2], L[3]),
		fmt.Sprintf(`'%s'"'%s''%s'\"\"'%s''%s'\''%s'`, L[0], L[1], L[2], L[3], L[4], L[5]),
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
