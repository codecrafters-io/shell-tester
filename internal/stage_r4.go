package internal

import (
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testR4(stageHarness *test_case_harness.TestCaseHarness) error {
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

	randomWords := random.RandomWords(3)
	slices.Sort(randomWords)

	filePaths := []string{
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[0])),
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[1])),
		path.Join(lsDir, fmt.Sprintf("%s", randomWords[2])),
	}
	fileContents := []string{
		randomWords[0] + "\n",
		randomWords[1] + "\n",
		randomWords[2] + "\n",
	}
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	randomWords2 := random.RandomWords(3)
	slices.Sort(randomWords2)
	outputFilePath := path.Join(stageDir, randomWords2[0]+".md")
	outputFilePath2 := path.Join(stageDir, randomWords2[1]+".md")
	// outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

	command1 := fmt.Sprintf("ls %s >> %s", "nonexistent", outputFilePath)

	reflectionTestCase := test_cases.CommandResponseTestCase{
		Command:          command1,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received error message",
	}

	asserter.AddAssertion(assertions.FileContentAssertion{
		FilePath:        outputFilePath,
		ExpectedContent: "",
	})

	if err := reflectionTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	command3 := fmt.Sprintf("ls %s 2>> %s", "nonexistent", outputFilePath2)
	command4 := fmt.Sprintf("cat %s", outputFilePath2)

	reflectionTestCase2 := test_cases.CommandReflectionTestCase{
		Command: command3,
	}
	if err := reflectionTestCase2.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command4,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
