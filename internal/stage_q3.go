package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f\\n%d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f\\%d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f'\\'%d", random.RandomInt(1, 100))),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	L := random.RandomElementsFromArray(LARGE_WORDS, 5)
	inputs := []string{
		`echo "before\   after"`,
		fmt.Sprintf(`echo %s\ \ \ \ \ \ %s`, L[0], L[1]),
		fmt.Sprintf(`echo %s\n%s`, L[2], L[3]),
		fmt.Sprintf(`cat "%s" "%s" "%s"`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		`before\   after`,
		fmt.Sprintf("%s      %s", L[0], L[1]),
		fmt.Sprintf("%sn%s", L[2], L[3]),
		fileContents[0] + fileContents[1] + strings.TrimRight(fileContents[2], "\n"),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
		testCase := test_cases.CommandResponseTestCase{
			Command:          testCaseContent.Input,
			ExpectedOutput:   testCaseContent.ExpectedOutput,
			FallbackPatterns: nil,
			SuccessMessage:   "Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	testCase := test_cases.CommandResponseTestCase{
		Command:        testCaseContents[3].Input,
		ExpectedOutput: testCaseContents[3].ExpectedOutput,
		SuccessMessage: "Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
