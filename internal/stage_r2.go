package internal

import (
	"fmt"
	"path"
	"slices"
	"strings"

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
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	dirs, err := getShortRandomDirectories(2)
	if err != nil {
		return err
	}
	stageDir := dirs[0]
	lsDir := dirs[1]
	defer cleanupDirectories(dirs)

	randomWords := random.RandomWords(1)
	slices.Sort(randomWords)

	filePaths := []string{
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[0])),
	}
	fileContents := []string{
		randomWords[0] + "\n",
	}
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	randomWords2 := random.RandomWords(3)
	slices.Sort(randomWords2)
	stringContent := strings.Join(randomWords, "\n")
	outputFilePath1 := path.Join(stageDir, randomWords2[0]+".md")
	outputFilePath2 := path.Join(stageDir, randomWords2[1]+".md")
	outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

	command1 := fmt.Sprintf("ls nonexistent 2> %s", outputFilePath1)
	command2 := fmt.Sprintf("cat %s", outputFilePath1)

	reflectionTestCase := test_cases.CommandReflectionTestCase{
		Command: command1,
	}
	if err := reflectionTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	responseTestCase1 := test_cases.CommandResponseTestCase{
		Command:          command2,
		ExpectedOutput:   "ls: nonexistent: No such file or directory",
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected error message",
	}

	if err := responseTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	stringContent = fmt.Sprintf("'%ss file not found'", getRandomName())
	command3 := fmt.Sprintf("echo %s 2> %s", stringContent, outputFilePath2)

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command3,
		ExpectedOutput:   stringContent[1 : len(stringContent)-1],
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
	logger.Successf("✓ File: %s passed content validation", outputFilePath2)

	///////

	file := filePaths[0]
	fileContent := randomWords[0]
	command5 := fmt.Sprintf("cat %s %s 2> %s", file, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("cat %s", outputFilePath3)

	responseTestCase2 := test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   fileContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received file content",
	}
	if err := responseTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	responseTestCase3 := test_cases.CommandResponseTestCase{
		Command:          command6,
		ExpectedOutput:   fmt.Sprintf("cat: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected error message",
	}

	if err := responseTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
