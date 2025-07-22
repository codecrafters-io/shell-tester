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

func testQ6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	executableName1 := `'exe  with  space'`
	executableName2 := `'exe with "quotes"'`
	executableName3 := `"exe with \'single quotes\'"`
	executableName4 := `'exe with \n newline'`
	executableDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "cat", CommandName: executableName1, CommandMetadata: ""},
		{CommandType: "cat", CommandName: executableName2, CommandMetadata: ""},
		{CommandType: "cat", CommandName: executableName3, CommandMetadata: ""},
		{CommandType: "cat", CommandName: executableName4, CommandMetadata: ""},
	}, true)
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	logAvailableExecutables(logger, []string{executableName1, executableName2, executableName3, executableName4})

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
	}

	if err := writeFiles([]string{
		path.Join(executableDir, "f1"),
		path.Join(executableDir, "f2"),
		path.Join(executableDir, "f3"),
		path.Join(executableDir, "f4"),
	}, []string{fileContents[0] + "\n", fileContents[1] + "\n", fileContents[2] + "\n", fileContents[3] + "\n"}, logger, nil); err != nil {
		return err
	}

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
