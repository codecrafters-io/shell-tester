package internal

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// TODO: Remove this function
//
//lint:ignore U1000 Ignore unused function warning for testAX
func testAX(stageHarness *test_case_harness.TestCaseHarness) error {
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

	command0 := fmt.Sprintf("typ" + "\t" + "\n")

	err = test_cases.CommandWithAttemptedCompletionTestCase{
		RawCommand:         command0,
		ExpectedReflection: "type",
		SuccessMessage:     "Received type",
		ExpectedAutocompletedReflectionHasNoSpace: true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	directory := "xyz"
	command1 := fmt.Sprintf("cd %s"+"\t"+"\n", directory[:1])

	err = test_cases.CommandWithAttemptedCompletionTestCase{
		RawCommand:         command1,
		ExpectedReflection: fmt.Sprintf("cd %s/", directory),
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SuccessMessage: fmt.Sprintf("Changed directory to %s", directory),
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	command2 := fmt.Sprintf("cat " + "\t" + "\n")
	err = test_cases.CommandWithAttemptedCompletionTestCase{
		RawCommand:         command2,
		ExpectedReflection: "cat file",
		ExpectedOutput:     "Hello World!",
		SuccessMessage:     "Received the contents of xyz/file",
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	randomDirectory, err := getRandomDirectory(stageHarness)
	if err != nil {
		return err
	}
	logger.Infof("Random directory: %s", randomDirectory)
	parentDirectory := path.Dir(randomDirectory)
	grandParentDirectory := path.Dir(parentDirectory)

	command3 := fmt.Sprintf("cd %s", grandParentDirectory)
	err = test_cases.CommandReflectionTestCase{
		Command:        command3,
		SuccessMessage: fmt.Sprintf("Changed directory to %s", grandParentDirectory),
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	command4 := "ls"
	err = test_cases.CommandResponseTestCase{
		Command:        command4,
		ExpectedOutput: filepath.Base(parentDirectory),
		SuccessMessage: fmt.Sprintf("Listed the contents of %s", grandParentDirectory),
	}.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	command5 := "cd " + "\t" + "\t" + "\n"
	autocompleted := fmt.Sprintf("cd %s/%s/", filepath.Base(parentDirectory), filepath.Base(randomDirectory))
	err = test_cases.CommandWithAttemptedCompletionTestCase{
		RawCommand:         command5,
		ExpectedReflection: autocompleted,
		ExpectedAutocompletedReflectionHasNoSpace: true,
		SuccessMessage: fmt.Sprintf("Changed directory to %s", filepath.Base(parentDirectory)),
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	command6 := "ec" + "\t" + "\t"
	err = test_cases.CommandWithAttemptedCompletionTestCase{
		RawCommand:         command6,
		ExpectedReflection: "ec",
		ExpectedAutocompletedReflectionHasNoSpace: true,
		ExpectedOutput:      "echo  ecpg",
		SuccessMessage:      "Received the contents of xyz/file",
		SkipPromptAssertion: true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
