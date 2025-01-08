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

func testR4(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the cls file is created)
	randomDir, err := getRandomDirectory()
	if err != nil {
		return err
	}

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
	// cls -1 nonexistent >> tmp.md

	customLsName := "cls"
	customLsPath := path.Join(randomDir, customLsName)
	err = custom_executable.CreateLsExecutable(customLsPath)
	if err != nil {
		return err
	}

	command1 := fmt.Sprintf("%s -1 %s >> %s", customLsName, "nonexistent", outputFilePath1)

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command1,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: []*regexp.Regexp{},
		SuccessMessage:   "✓ Received error message",
	}
	asserter.AddAssertion(assertions.FileContentAssertion{
		FilePath:        outputFilePath1,
		ExpectedContent: "",
	})
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}
	logger.Successf("✓ File: %s is empty", outputFilePath1)

	// Test2:
	// cls -1 nonexistent 2>> tmp.md

	command2 := fmt.Sprintf("%s -1 %s 2>> %s", customLsName, "nonexistent", outputFilePath2)
	command3 := fmt.Sprintf("cat %s", outputFilePath2)

	err = test_cases.CommandReflectionTestCase{
		Command: command2,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command3,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: []*regexp.Regexp{},
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test3:
	// echo "Error" 2>> tmp.md
	// cat nonexistent 2>> tmp.md
	// cls -1 nonexistent 2>> tmp.md
	// cat tmp.md

	message := fmt.Sprintf("%s says Error", getRandomName())
	command4 := fmt.Sprintf(`echo "%s" 2>> %s`, message, outputFilePath3)
	command5 := fmt.Sprintf(`cat %s 2>> %s`, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("%s -1 %s 2>> %s", customLsName, "nonexistent", outputFilePath3)
	command7 := fmt.Sprintf("cat %s", outputFilePath3)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command4,
		ExpectedOutput:   message,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}
	err = test_cases.CommandReflectionTestCase{
		Command: command5,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	err = test_cases.CommandReflectionTestCase{
		Command: command6,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	errorMessagesInFile := []string{
		"cat: nonexistent: No such file or directory",
		"ls: nonexistent: No such file or directory",
	}

	linuxLSErrorMessage := errorMessagesInFile[1]
	linuxLSErrorMessageRegex := []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf("^%s$", linuxLSErrorMessage))}
	alpineCatErrorMessage := "cat: can't open 'nonexistent': No such file or directory"
	alpineCatErrorMessageRegex := []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf("^%s$", alpineCatErrorMessage))}

	// TODO: Simplify this after writing custom cat
	multiLineAssertion := assertions.NewEmptyMultiLineAssertion()
	multiLineAssertion.AddSingleLineAssertion(errorMessagesInFile[0], alpineCatErrorMessageRegex)
	multiLineAssertion.AddSingleLineAssertion(errorMessagesInFile[1], linuxLSErrorMessageRegex)

	multiLineResponseTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:            command7,
		MultiLineAssertion: multiLineAssertion,
		SuccessMessage:     "✓ Received redirected file content",
	}
	if err := multiLineResponseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
