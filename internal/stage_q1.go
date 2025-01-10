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

func testQ1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []string{"cat"})
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := getShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	randomUniqueFileNames := random.RandomInts(1, 100, 3)
	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f   %d", randomUniqueFileNames[0])),
		path.Join(randomDir, fmt.Sprintf("f   %d", randomUniqueFileNames[1])),
		path.Join(randomDir, fmt.Sprintf("f   %d", randomUniqueFileNames[2])),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	L := random.RandomElementsFromArray(LARGE_WORDS, 5)
	inputs := []string{
		fmt.Sprintf(`echo '%s %s'`, L[0], L[1]),
		fmt.Sprintf(`echo %s     %s`, L[1], L[4]),
		fmt.Sprintf(`echo '%s     %s' '%s''%s'`, L[2], L[3], L[4], L[0]),
		fmt.Sprintf(`%s '%s' '%s' '%s'`, CUSTOM_CAT_COMMAND, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf("%s %s", L[0], L[1]),
		fmt.Sprintf("%s %s", L[1], L[4]),
		fmt.Sprintf("%s     %s %s%s", L[2], L[3], L[4], L[0]),
		fileContents[0] + fileContents[1] + strings.TrimRight(fileContents[2], "\n"),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
		testCase := test_cases.CommandResponseTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "✓ Received expected response",
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
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

type testCaseContent struct {
	Input          string
	ExpectedOutput string
}

func newTestCaseContent(input string, expectedOutput string) testCaseContent {
	return testCaseContent{
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

func newTestCaseContents(inputs []string, expectedOutputs []string) []testCaseContent {
	testCases := []testCaseContent{}
	for i, input := range inputs {
		testCases = append(testCases, newTestCaseContent(input, expectedOutputs[i]))
	}
	return testCases
}
