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

func testQ5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("\"f %d\"", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("\"f\\%d\"", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f%d", random.RandomInt(1, 100))),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}

	if err := shell.Start(); err != nil {
		return err
	}

	_, L := getRandomWordsSmallAndLarge(5, 5)
	inputs := []string{
		fmt.Sprintf(`echo "%s'%s'\\n'%s"`, L[0], L[1], L[2]),
		fmt.Sprintf(`echo "%s\"insidequotes"%s\"`, L[0], L[1]),
		fmt.Sprintf(`echo "mixed\"quote'%s'\\"`, L[4]),
		fmt.Sprintf(`cat '%s' '%s' '%s'`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf(`%s'%s'\n'%s`, L[0], L[1], L[2]),
		fmt.Sprintf(`%s"insidequotes%s"`, L[0], L[1]),
		fmt.Sprintf(`mixed"quote'%s'\`, L[4]),
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
