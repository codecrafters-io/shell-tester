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

func testR2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []string{"ls", "cat"})
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	dirs, err := getShortRandomDirectories(stageHarness, 2)
	if err != nil {
		return err
	}
	stageDir, lsDir := dirs[0], dirs[1]

	randomWords := random.RandomWords(1)
	slices.Sort(randomWords)
	filePaths := []string{
		path.Join(lsDir, randomWords[0]),
	}
	fileContents := []string{
		randomWords[0] + "\n",
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
	// ls -1 nonexistent 2> tmp.md; cat tmp.md

	command1 := fmt.Sprintf("%s -1 nonexistent 2> %s", CUSTOM_LS_COMMAND, outputFilePath1)
	command2 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath1)

	err = test_cases.CommandReflectionTestCase{
		Command: command1,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command2,
		ExpectedOutput:   "ls: nonexistent: No such file or directory",
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected error message",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test2:
	// echo 'File not found' 2> tmp.md; cat tmp.md

	message := fmt.Sprintf("%s file cannot be found", getRandomName())
	command3 := fmt.Sprintf("echo %s 2> %s", fmt.Sprintf("'%s'", message), outputFilePath2)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command3,
		ExpectedOutput:   message,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected error message",
	}
	asserter.AddAssertion(assertions.FileContentAssertion{
		FilePath:        outputFilePath2,
		ExpectedContent: "",
	})
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}
	logger.Successf("✓ File: %s is empty", outputFilePath2)

	// Test3:
	// cat exists nonexistent 2> tmp.md; cat tmp.md

	filePath := filePaths[0]
	fileContent := randomWords[0]
	command5 := fmt.Sprintf("%s %s %s 2> %s", CUSTOM_CAT_COMMAND, filePath, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath3)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   fileContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command6,
		ExpectedOutput:   fmt.Sprintf("cat: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected error message",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
