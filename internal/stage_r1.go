package internal

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testR1(stageHarness *test_case_harness.TestCaseHarness) error {
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

	randomWords := random.RandomWords(2)
	slices.Sort(randomWords)

	filePaths := []string{
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[0])),
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[1])),
	}
	fileContents := []string{
		randomWords[0] + "\n",
		randomWords[1] + "\n",
	}
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	stringContent := strings.Join(randomWords, "\n")
	outputFilePath := path.Join(stageDir, random.RandomWord()+".md")
	outputFilePath2 := path.Join(stageDir, random.RandomWord()+".md")
	outputFilePath3 := path.Join(stageDir, random.RandomWord()+".md")
	command3 := fmt.Sprintf("ls %s > %s", lsDir, outputFilePath)
	command4 := fmt.Sprintf("cat %s", outputFilePath)

	reflectionTestCase := test_cases.CommandReflectionTestCase{
		Command: command3,
	}
	if err := reflectionTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	responseTestCase1 := test_cases.CommandWithMultilineResponseTestCase{
		Command:          command4,
		ExpectedOutput:   randomWords,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase1.Run(asserter, shell, logger); err != nil {
		return err
	}

	stringContent = "Hello " + getRandomName()
	command1 := fmt.Sprintf("echo '%s' 1> %s", stringContent, outputFilePath2)
	command2 := fmt.Sprintf("cat %s", outputFilePath2)

	reflectionTestCase2 := test_cases.CommandReflectionTestCase{
		Command: command1,
	}
	if err := reflectionTestCase2.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command2,
		ExpectedOutput:   stringContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	file := filePaths[1]
	fileContent := randomWords[1]
	command5 := fmt.Sprintf("cat %s %s 1> %s", file, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("cat %s", outputFilePath3)

	reflectionTestCase3 := test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   fmt.Sprintf("cat: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received error message",
	}
	if err := reflectionTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	responseTestCase3 := test_cases.CommandResponseTestCase{
		Command:          command6,
		ExpectedOutput:   fileContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
