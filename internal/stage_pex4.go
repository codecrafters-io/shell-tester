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

	words := random.RandomWords(6)

	// Valid assignment: name starting with an underscore
	validVariable := "_" + words[0]
	validVariableValue := words[1]
	underscorePrefixedAssignment := test_cases.DeclareAssignmentTestCase{Variable: validVariable, Value: validVariableValue}
	if err := underscorePrefixedAssignment.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Invalid assignment: name starting with a digit
	invalidDigitPrefixedVariable := fmt.Sprintf("%d%s", random.RandomInt(2, 9), words[2])
	invalidVariableValue := words[3]
	digitPrefixedAssignment := test_cases.DeclareAssignmentTestCase{Variable: invalidDigitPrefixedVariable, Value: invalidVariableValue}
	if err := digitPrefixedAssignment.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Invalid assignment: name with a hyphen in the middle
	hyphenInMiddleVariable := words[4] + "-" + words[5]
	hyphenInMiddleAssignment := test_cases.DeclareAssignmentTestCase{Variable: hyphenInMiddleVariable, Value: validVariableValue}
	if err := hyphenInMiddleAssignment.Run(asserter, shell, logger); err != nil {
		return err
	}

	printValidVariable := test_cases.DeclarePrintTestCase{Variable: validVariable, Value: validVariableValue}
	if err := printValidVariable.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
