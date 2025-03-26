package internal

import (
	"fmt"
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
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	randomUniqueFileNames := random.RandomInts(1, 100, 3)
	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f\\n%d", randomUniqueFileNames[0])),
		path.Join(randomDir, fmt.Sprintf("f\\%d", randomUniqueFileNames[1])),
		path.Join(randomDir, fmt.Sprintf("f'\\'%d", randomUniqueFileNames[2])),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	L := random.RandomElementsFromArray(LARGE_WORDS, 5)
	inputs := []string{
		fmt.Sprintf(`echo %s\ \ \ \ \ \ %s`, L[0], L[1]),
		fmt.Sprintf(`echo \'\"%s %s\"\'`, L[1], L[2]),
		fmt.Sprintf(`echo %s\n%s`, L[2], L[3]),
		fmt.Sprintf(`%s "%s" "%s" "%s"`, CUSTOM_CAT_COMMAND, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf("%s      %s", L[0], L[1]),
		fmt.Sprintf(`'"%s %s"'`, L[1], L[2]),
		fmt.Sprintf("%sn%s", L[2], L[3]),
		fileContents[0] + fileContents[1] + strings.TrimRight(fileContents[2], "\n"),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
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

	testCase := test_cases.CommandResponseTestCase{
		Command:        testCaseContents[3].Input,
		ExpectedOutput: testCaseContents[3].ExpectedOutput,
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
