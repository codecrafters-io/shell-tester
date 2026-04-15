package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPEX3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	words := random.RandomWords(2)
	variableName := words[0]
	variableValue := words[1]

	assignTestCase := test_cases.DeclareAssignmentTestCase{
		Variable: variableName,
		Value:    variableValue,
	}
	if err := assignTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	printTestCase := test_cases.DeclarePrintTestCase{
		Variable: variableName,
		Value:    variableValue,
	}
	if err := printTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	missingVariable := "missing_" + random.RandomWords(1)[0]
	printErrorTestCase := test_cases.DeclarePrintErrorTestCase{
		Variable: missingVariable,
	}
	if err := printErrorTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
