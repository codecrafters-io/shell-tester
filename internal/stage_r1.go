package internal

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testR1(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the cls file is created)
	randomDir, err := getRandomDirectory()
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(randomDir)
	}()

	pathEnvVar := os.Getenv("PATH")
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, pathEnvVar))
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	dirs, err := getShortRandomDirectories(2)
	if err != nil {
		return err
	}
	stageDir, lsDir := dirs[0], dirs[1]
	defer cleanupDirectories(dirs)

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
	// cls -1 foo > tmp.md; cat tmp.md

	customLsName := "cls"
	customLsPath := path.Join(randomDir, customLsName)
	err = custom_executable.CreateLsExecutable(customLsPath)
	if err != nil {
		return err
	}

	command1 := fmt.Sprintf("%s -1 %s > %s", customLsName, lsDir, outputFilePath1)
	command2 := fmt.Sprintf("cat %s", outputFilePath1)

	err = test_cases.CommandReflectionTestCase{
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
	command4 := fmt.Sprintf("cat %s", outputFilePath2)

	err = test_cases.CommandReflectionTestCase{
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
	command5 := fmt.Sprintf("cat %s %s 1> %s", filePath, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("cat %s", outputFilePath3)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command5,
		ExpectedOutput:   fmt.Sprintf("cat: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf("cat: can't open '%s': No such file or directory", "nonexistent"))},
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

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
