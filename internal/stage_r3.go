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

func testR3(stageHarness *test_case_harness.TestCaseHarness) error {
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
	outputFilePath := path.Join(stageDir, randomWords2[0]+".md")
	outputFilePath2 := path.Join(stageDir, randomWords2[1]+".md")
	outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

	// Test1:
	// ls -1 foo >> tmp.md; cat tmp.md

	command1 := fmt.Sprintf("%s -1 %s >> %s", CUSTOM_LS_COMMAND, lsDir, outputFilePath)
	command2 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command1,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	responseTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:            command2,
		MultiLineAssertion: assertions.NewMultiLineAssertion(randomWords),
		SuccessMessage:     "✓ Received redirected file content",
	}

	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test2:
	// echo 'Hello' 1>> tmp.md; echo 'Hello' 1>> tmp.md; cat tmp.md

	message1 := "Hello " + getRandomName()
	message2 := "Hello " + getRandomName()
	command4 := fmt.Sprintf("echo '%s' 1>> %s", message1, outputFilePath2)
	command5 := fmt.Sprintf("echo '%s' 1>> %s", message2, outputFilePath2)
	command6 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath2)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command4,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command5,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	responseTestCase = test_cases.CommandWithMultilineResponseTestCase{
		Command:            command6,
		MultiLineAssertion: assertions.NewMultiLineAssertion([]string{message1, message2}),
		SuccessMessage:     "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test3:
	// echo "List of files: " > tmp.md; ls -1 foo >> tmp.md; cat tmp.md

	command7 := fmt.Sprintf(`echo "List of files: " > %s`, outputFilePath3)
	command8 := fmt.Sprintf("%s -1 %s >> %s", CUSTOM_LS_COMMAND, lsDir, outputFilePath3)
	command9 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath3)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command7,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command8,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}

	responseTestCase = test_cases.CommandWithMultilineResponseTestCase{
		Command:            command9,
		MultiLineAssertion: assertions.NewMultiLineAssertion(append([]string{"List of files:"}, randomWords...)),
		SuccessMessage:     "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
