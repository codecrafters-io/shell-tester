package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f %d", getUniqueRandomIntegerFileNames(1, 100, 3)[0])),
		path.Join(randomDir, fmt.Sprintf("f   %d", getUniqueRandomIntegerFileNames(1, 100, 3)[1])),
		path.Join(randomDir, fmt.Sprintf("f's%d", getUniqueRandomIntegerFileNames(1, 100, 3)[2])),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}

	if err := shell.Start(); err != nil {
		return err
	}

	L := random.RandomElementsFromArray(LARGE_WORDS, 5)
	inputs := []string{
		fmt.Sprintf(`echo "%s %s"`, L[0], L[1]),
		fmt.Sprintf(`echo "%s  %s"  "%s"`, L[1], L[2], L[3]),
		fmt.Sprintf(`echo "%s"  "%s's"  "%s"`, L[3], L[4], L[1]),
		fmt.Sprintf(`cat "%s" "%s" "%s"`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf("%s %s", L[0], L[1]),
		fmt.Sprintf("%s  %s %s", L[1], L[2], L[3]),
		fmt.Sprintf(`%s %s's %s`, L[3], L[4], L[1]),
		fileContents[0] + fileContents[1] + strings.TrimRight(fileContents[2], "\n"),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}
	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:        testCaseContents[3].Input,
		ExpectedOutput: testCaseContents[3].ExpectedOutput,
		SuccessMessage: "Received expected response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
