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

func testP6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
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
	randomWords := random.RandomWords(5)
	fileContent := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4])
	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	lines := strings.Count(fileContent, "\n")
	words := strings.Count(strings.ReplaceAll(fileContent, "\n", " "), " ") + 1
	bytes := len(fileContent)

	input := fmt.Sprintf(`cat %s | head -n 5 | wc`, filePath)
	expectedOutput := fmt.Sprintf("%7d%8d%8d", lines, words, bytes)

	singleLineTestCase := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
