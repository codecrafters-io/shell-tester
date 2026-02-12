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

func testR4(stageHarness *test_case_harness.TestCaseHarness) error {
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

	dirs, err := CreateShortRandomDirsInTmp(stageHarness, 1)
	if err != nil {
		return err
	}
	stageDir := dirs[0]

	randomWords2 := random.RandomElementsFromArray(SMALL_WORDS, 3)
	slices.Sort(randomWords2)
	outputFilePath1 := path.Join(stageDir, randomWords2[0]+".md")
	outputFilePath2 := path.Join(stageDir, randomWords2[1]+".md")
	outputFilePath3 := path.Join(stageDir, randomWords2[2]+".md")

	// Test1:
	// ls -1 nonexistent >> tmp.md

	command1 := fmt.Sprintf("%s -1 %s >> %s", CUSTOM_LS_COMMAND, "nonexistent", outputFilePath1)

	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command1,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received error message",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Assert file is empty after the command response test case
	// to avoid printing the error message before the command is reflected
	asserter.AddAssertion(assertions.FileContentAssertion{
		FilePath:        outputFilePath1,
		ExpectedContent: "",
	})
	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("✓ File: %s is empty", outputFilePath1)

	// Test2:
	// ls -1 nonexistent 2>> tmp.md

	command2 := fmt.Sprintf("%s -1 %s 2>> %s", CUSTOM_LS_COMMAND, "nonexistent", outputFilePath2)
	command3 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath2)

	err = test_cases.CommandWithNoResponseTestCase{
		Command: command2,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command3,
		ExpectedOutput:   fmt.Sprintf("ls: %s: No such file or directory", "nonexistent"),
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test3:
	// echo "Error" 2>> tmp.md
	// cat nonexistent 2>> tmp.md
	// ls -1 nonexistent 2>> tmp.md
	// cat tmp.md

	message := fmt.Sprintf("%s says Error", getRandomName())
	command4 := fmt.Sprintf(`echo "%s" 2>> %s`, message, outputFilePath3)
	command5 := fmt.Sprintf(`%s %s 2>> %s`, CUSTOM_CAT_COMMAND, "nonexistent", outputFilePath3)
	command6 := fmt.Sprintf("%s -1 %s 2>> %s", CUSTOM_LS_COMMAND, "nonexistent", outputFilePath3)
	command7 := fmt.Sprintf("%s %s", CUSTOM_CAT_COMMAND, outputFilePath3)

	responseTestCase = test_cases.CommandResponseTestCase{
		Command:          command4,
		ExpectedOutput:   message,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}
	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}
	err = test_cases.CommandWithNoResponseTestCase{
		Command: command5,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	err = test_cases.CommandWithNoResponseTestCase{
		Command: command6,
	}.Run(asserter, shell, logger, true)
	if err != nil {
		return err
	}
	errorMessagesInFile := []string{
		"cat: nonexistent: No such file or directory",
		"ls: nonexistent: No such file or directory",
	}

	multiLineResponseTestCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:            command7,
		MultiLineAssertion: assertions.NewMultiLineAssertion(errorMessagesInFile),
		SuccessMessage:     "✓ Received redirected file content",
	}
	if err := multiLineResponseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
