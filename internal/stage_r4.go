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
	outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

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

	// //////

	message := fmt.Sprintf("%s says Error", getRandomName())
	command5_1 := fmt.Sprintf(`echo "%s" 2>> %s`, message, outputFilePath3)
	command5_2 := fmt.Sprintf(`cat %s 2>> %s`, "nonexistent", outputFilePath3)
	command5_3 := fmt.Sprintf("ls %s 2>> %s", "nonexistent", outputFilePath3)
	command5_4 := fmt.Sprintf("cat %s", outputFilePath3)

	responseTestCase5_1 := test_cases.CommandResponseTestCase{
		Command:          command5_1,
		ExpectedOutput:   message,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase5_1.Run(asserter, shell, logger); err != nil {
		return err
	}

	reflectionTestCase5_2 := test_cases.CommandReflectionTestCase{
		Command: command5_2,
	}
	if err := reflectionTestCase5_2.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	reflectionTestCase5_3 := test_cases.CommandReflectionTestCase{
		Command: command5_3,
	}
	if err := reflectionTestCase5_3.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	errorMessagesInFile := []string{
		"cat: nonexistent: No such file or directory",
		"ls: nonexistent: No such file or directory",
	}
	responseTestCase5_4 := test_cases.CommandWithMultilineResponseTestCase{
		Command:          command5_4,
		ExpectedOutput:   errorMessagesInFile,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase5_4.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
