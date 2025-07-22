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

func testP1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
		{CommandType: "tail", CommandName: CUSTOM_TAIL_COMMAND, CommandMetadata: ""},
		{CommandType: "wc", CommandName: CUSTOM_WC_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Test-1
	randomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	filePath := path.Join(randomDir, fmt.Sprintf("file-%d", random.RandomInt(1, 100)))
	randomWords := random.RandomWords(10)
	fileContent := fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4], randomWords[5], randomWords[6], randomWords[7], randomWords[8], randomWords[9])

	lines := strings.Count(fileContent, "\n")
	words := strings.Count(strings.ReplaceAll(fileContent, "\n", " "), " ")
	bytes := len(fileContent)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	input := fmt.Sprintf(`cat %s | wc`, filePath)
	expectedOutput := fmt.Sprintf("%8d%8d%8d", lines, words, bytes)

	if err := writeFiles([]string{filePath}, []string{fileContent}, logger, &WriteFilesOptions{EchoWithoutFlagN: true}); err != nil {
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

	// Test-2
	randomDir, err = GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	filePath = path.Join(randomDir, fmt.Sprintf("file-%d", random.RandomInt(1, 100)))
	randomWords = random.RandomWords(6)
	fileContent = fmt.Sprintf("1. %s %s\n2. %s %s\n3. %s %s\n", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4], randomWords[5])
	if err := writeFiles([]string{filePath}, []string{fileContent}, logger, &WriteFilesOptions{EchoWithoutFlagN: true}); err != nil {
		return err
	}

	input = fmt.Sprintf(`tail -f %s | head -n 5`, filePath)
	expectedMultiLineOutput := strings.Split(strings.Trim(fileContent, "\n"), "\n")
	multiLineTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:             input,
		MultiLineAssertion:  assertions.NewMultiLineAssertion(expectedMultiLineOutput),
		SuccessMessage:      "✓ Received redirected file content",
		SkipPromptAssertion: true,
	}
	if err := multiLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Append content to the file while command is running
	line4 := "4. " + strings.Join(random.RandomWords(2), " ")
	if err := appendFile(filePath, fmt.Sprintf("%s\n", line4)); err != nil {
		return err
	}

	firstSingleLineAssertion := assertions.SingleLineAssertion{
		ExpectedOutput: line4,
	}
	asserter.AddAssertion(&firstSingleLineAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	logger.Successf("✓ Received appended line 4")

	// Append again
	line5 := "5. " + strings.Join(random.RandomWords(2), " ")
	if err := appendFile(filePath, fmt.Sprintf("%s\n", line5)); err != nil {
		return err
	}

	secondSingleLineAssertion := assertions.SingleLineAssertion{
		ExpectedOutput: line5,
	}
	asserter.AddAssertion(&secondSingleLineAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}
	logger.Successf("✓ Received appended line 5")

	return logAndQuit(asserter, nil)
}
