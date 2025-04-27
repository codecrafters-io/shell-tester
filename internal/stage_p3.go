package internal

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testP3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
		{CommandType: "tail", CommandName: CUSTOM_TAIL_COMMAND, CommandMetadata: ""},
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
	randomWords := random.RandomWords(6)
	fileContent := fmt.Sprintf("%s %s\n%s %s\n%s %s\n", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4], randomWords[5])

	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	input := fmt.Sprintf(`tail -f %s | head -n 5`, filePath)
	expectedOutput := strings.Split(fileContent, "\n")
	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	multiLineTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:             input,
		MultiLineAssertion:  assertions.NewMultiLineAssertion(expectedOutput),
		SuccessMessage:      "✓ Received redirected file content",
		SkipPromptAssertion: true,
	}
	if err := multiLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Append content to the file while command is running
	if err := appendFile(filePath, "This is line 4.\n"); err != nil {
		return err
	}

	firstSingleLineAssertion := assertions.SingleLineAssertion{
		ExpectedOutput: "This is line 4.",
	}
	asserter.AddAssertion(&firstSingleLineAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	logger.Successf("✓ Received appended line 4")

	// Append again
	if err := appendFile(filePath, "This is line 5.\n"); err != nil {
		return err
	}

	secondSingleLineAssertion := assertions.SingleLineAssertion{
		ExpectedOutput: "This is line 5.",
	}
	asserter.AddAssertion(&secondSingleLineAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	logger.Successf("✓ Received appended line 5")

	return logAndQuit(asserter, nil)
}
