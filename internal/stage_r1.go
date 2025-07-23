package internal

import (
	"fmt"
	"path"
	"slices"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testR1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "ls", CommandName: CUSTOM_LS_COMMAND, CommandMetadata: ""},
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	dirs, err := GetShortRandomDirectories(stageHarness, 2)
	if err != nil {
		return err
	}
	stageDir, lsDir := dirs[0], dirs[1]

	randomWords := random.RandomWords(3)
	slices.Sort(randomWords)
	filePaths := []string{
		path.Join(lsDir, randomWords[0]),
		path.Join(lsDir, randomWords[1]),
		path.Join(lsDir, randomWords[2]),
	}
	fileContents := []string{
		randomWords[0] + "\n",
		randomWords[1] + "\n",
		randomWords[2] + "\n",
	}
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	randomWords2 := random.RandomElementsFromArray(SMALL_WORDS, 3)
	slices.Sort(randomWords2)
	outputFilePath1 := path.Join(stageDir, randomWords2[0]+".md")
	outputFilePath2 := path.Join(stageDir, randomWords2[1]+".md")
	outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

	// Test1:
	// ls -1 foo > tmp.md; cat tmp.md

	command1 := fmt.Sprintf("%s -1 %s > %s", CUSTOM_LS_COMMAND, lsDir, outputFilePath1)
	command2 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath1)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command1,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	multiLineTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:            command2,
		MultiLineAssertion: assertions.NewMultiLineAssertion(randomWords),
		SuccessMessage:     "✓ Received redirected file content",
	}
	if err := multiLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test2:
	// echo "Hello Ryan" 1> tmp.md; cat tmp.md

	message := "Hello " + getRandomName()
	command3 := fmt.Sprintf("echo '%s' 1> %s", message, outputFilePath2)
	command4 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath2)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command3,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command4,
		ExpectedOutput:   message,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test3:
	// cat exists nonexistent > tmp.md; cat tmp.md

	filePath := filePaths[1]
	fileContent := randomWords[1]
	command5 := fmt.Sprintf("%s %s %s 1> %s", CUSTOM_CAT_COMMAND, filePath, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath3)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   fmt.Sprintf("cat: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received error message",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command6,
		ExpectedOutput:   fileContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
