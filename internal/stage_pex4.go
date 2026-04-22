package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPEX4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	words := random.RandomWords(4)

	// Valid assignment: underscore-prefixed name
	validVariable := "_" + words[0]
	validVariableValue := words[1]
	if err := (test_cases.DeclareAssignmentTestCase{Variable: validVariable, Value: validVariableValue}).Run(asserter, shell, logger); err != nil {
		return err
	}

	// Invalid assignment: name starting with a digit
	invalidVariable := fmt.Sprintf("%d%s", random.RandomInt(2, 9), words[2])
	invalidVariableValue := words[3]
	if err := (test_cases.DeclareAssignmentTestCase{Variable: invalidVariable, Value: invalidVariableValue}).Run(asserter, shell, logger); err != nil {
		return err
	}

	if err := (test_cases.DeclarePrintTestCase{Variable: validVariable, Value: validVariableValue}).Run(asserter, shell, logger); err != nil {
		return err
	}

	if err := (test_cases.DeclarePrintErrorTestCase{Variable: invalidVariable}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
