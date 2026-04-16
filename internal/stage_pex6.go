package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPEX6(stageHarness *test_case_harness.TestCaseHarness) error {
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

	// Declare two variables with valid names and random values
	words := random.RandomWords(6)
	integerSuffixes := random.RandomInts(1, 10, 2)
	variableName1 := fmt.Sprintf(
		"%s_%d",
		strings.ToUpper(words[0][:1])+words[0][1:],
		integerSuffixes[0],
	)
	variableName2 := fmt.Sprintf(
		"%s_%d",
		strings.ToUpper(words[1][:1])+words[1][1:],
		integerSuffixes[1],
	)
	variableValue1 := words[2]
	variableValue2 := words[3]
	literalPrefix := words[4]
	literalSuffix := words[5]

	if err := (test_cases.DeclareAssignmentTestCase{Variable: variableName1, Value: variableValue1}).Run(asserter, shell, logger); err != nil {
		return err
	}
	if err := (test_cases.DeclareAssignmentTestCase{Variable: variableName2, Value: variableValue2}).Run(asserter, shell, logger); err != nil {
		return err
	}

	// Run the executable with ${VAR1} embedded between random literal prefix/suffix,
	// and ${VAR2} as a standalone brace-expansion argument.
	argument1 := fmt.Sprintf("%s_${%s}_%s", literalPrefix, variableName1, literalSuffix)
	argument2 := fmt.Sprintf("${%s}", variableName2)
	command := fmt.Sprintf("%s %s %s", executableName, argument1, argument2)

	expectedArgument1 := fmt.Sprintf("%s_%s_%s", literalPrefix, variableValue1, literalSuffix)
	expectedArgument2 := variableValue2

	expectedLines := []string{
		"Program was passed 3 args (including program name).",
		fmt.Sprintf("Arg #0 (program name): %s", executableName),
		fmt.Sprintf("Arg #1: %s", expectedArgument1),
		fmt.Sprintf("Arg #2: %s", expectedArgument2),
		fmt.Sprintf("Program Signature: %s", commandMetadata),
	}

	testCase := test_cases.CommandWithMultilineResponseTestCase{
		Command:            strings.Join([]string{command}, " "),
		MultiLineAssertion: assertions.NewMultiLineAssertion(expectedLines),
		SuccessMessage:     "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
