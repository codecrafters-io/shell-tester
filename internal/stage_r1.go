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
	for _, dir := range dirs {
		defer os.RemoveAll(dir)
	}

	randomWords := random.RandomWords(2)
	filePaths := []string{
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[0])),
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[1])),
	}
	if err := writeFiles(filePaths, randomWords, logger); err != nil {
		return err
	}

	slices.Sort(randomWords)
	stringContent := strings.Join(randomWords, "\n")
	outputFilePath := path.Join(stageDir, random.RandomWord()+".md")
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
	outputFilePath = path.Join(stageDir, random.RandomWord()+".md")
	command1 := fmt.Sprintf("echo '%s' 1> %s", stringContent, outputFilePath)
	command2 := fmt.Sprintf("cat %s", outputFilePath)

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

	outputFilePath2 := path.Join(stageDir, random.RandomWord()+".md")
	command5 := fmt.Sprintf("cat %s %s 1> %s", outputFilePath, "nonexistent", outputFilePath2)
	command6 := fmt.Sprintf("cat %s", outputFilePath2)

	reflectionTestCase3 := test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   "cat: nonexistent: No such file or directory",
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received error message",
	}
	if err := reflectionTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	responseTestCase3 := test_cases.CommandResponseTestCase{
		Command:          command6,
		ExpectedOutput:   stringContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase3.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
