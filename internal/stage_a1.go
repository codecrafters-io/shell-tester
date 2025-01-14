package internal

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testA1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []string{"ls", "cat"})
	if err != nil {
		return err
	}
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
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

func logScreenState(shell *shell_executable.ShellExecutable) {
	screenState := shell.GetScreenState()
	for _, row := range screenState {
		cleanedRow := virtual_terminal.BuildCleanedRow(row)
		if len(cleanedRow) > 0 {
			fmt.Println(cleanedRow)
		}
	}
}
