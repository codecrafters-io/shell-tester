package internal

import (
	"fmt"
	"os"
	"path"

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

	stageDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(stageDir)

	stringContent := "Hello " + getRandomName()
	outputFilePath := path.Join(stageDir, random.RandomWord()+".md")
	reflectionTestCase := test_cases.CommandReflectionTestCase{
		Command: fmt.Sprintf("echo '%s' 1> %s", stringContent, outputFilePath),
	}
	if err := reflectionTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	command := fmt.Sprintf("cat %s", outputFilePath)
	responseTestCase := test_cases.CommandResponseTestCase{
		Command:          command,
		ExpectedOutput:   stringContent,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received redirected file content",
	}

	if err := responseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

// Possibly catTestCase
// EchoTestCase
// LSTestCase

// FileContentTestCase
// AppendedFileContentTestCase
