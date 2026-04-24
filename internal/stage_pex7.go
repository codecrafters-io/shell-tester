package internal

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPEX7(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	commandMetadata := getRandomString()
	executableName := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))

	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executableName, CommandMetadata: commandMetadata},
	}, true)
	if err != nil {
		return err
	}
	logAvailableExecutables(logger, []string{executableName})

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	words := random.RandomWords(2)
	variableName := words[0]
	variableValue := words[1]

	assignVariable := test_cases.DeclareAssignmentTestCase{Variable: variableName, Value: variableValue}
	if err := assignVariable.Run(asserter, shell, logger); err != nil {
		return err
	}

	integerSuffixes := random.RandomInts(0, 10, 3)

	command := fmt.Sprintf(
		"%s ${missing_var_%d}_suffix ${%s} ${missing_var_%d} $missing_var_%d",
		executableName,
		integerSuffixes[0],
		variableName,
		integerSuffixes[1],
		integerSuffixes[2],
	)

	expectedLines := []string{
		"Program was passed 3 args (including program name).",
		fmt.Sprintf("Arg #0 (program name): %s", executableName),
		"Arg #1: _suffix",
		fmt.Sprintf("Arg #2: %s", variableValue),
		fmt.Sprintf("Program Signature: %s", commandMetadata),
	}

	missingVarExpansionCall := test_cases.CommandWithMultilineResponseTestCase{
		Command:            command,
		MultiLineAssertion: assertions.NewMultiLineAssertion(expectedLines),
		SuccessMessage:     "✓ Received expected response",
	}

	if err := missingVarExpansionCall.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
