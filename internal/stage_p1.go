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

func testP1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "wc", CommandName: CUSTOM_WC_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	filePath := path.Join(randomDir, fmt.Sprintf("file-%d", random.RandomInt(1, 100)))
	randomWords := random.RandomWords(10)
	fileContent := fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %s\n%s %s", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4], randomWords[5], randomWords[6], randomWords[7], randomWords[8], randomWords[9])

	lines := strings.Count(fileContent, "\n")
	words := strings.Count(strings.ReplaceAll(fileContent, "\n", " "), " ") + 1
	bytes := len(fileContent)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	input := fmt.Sprintf(`cat %s | wc`, filePath)
	expectedOutput := fmt.Sprintf("%7d%8d%8d", lines, words, bytes)

	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	testCase := test_cases.CommandResponseTestCase{
		Command:        input,
		ExpectedOutput: expectedOutput,
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
