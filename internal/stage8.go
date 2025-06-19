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

func testRun(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomCode1 := getRandomString()
	randomCode2 := getRandomString()

	randomName1 := getRandomName()
	randomName2 := getRandomName()
	randomName3 := getRandomName()

	randomExecutableName1 := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))
	randomExecutableName2 := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))

	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: randomExecutableName1, CommandMetadata: randomCode1},
		{CommandType: "signature_printer", CommandName: randomExecutableName2, CommandMetadata: randomCode2},
	}, true)
	if err != nil {
		return err
	}
	logAvailableExecutables(logger, []string{randomExecutableName1, randomExecutableName2})
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	command1 := []string{
		randomExecutableName1, randomName1,
	}

	testCase := test_cases.CommandWithMultilineResponseTestCase{
		Command: strings.Join(command1, " "),
		MultiLineAssertion: assertions.NewMultiLineAssertion([]string{
			fmt.Sprintf("Program was passed %d args (including program name).", len(command1)),
			fmt.Sprintf("Arg #0 (program name): %s", command1[0]),
			fmt.Sprintf("Arg #1: %s", command1[1]),
			fmt.Sprintf("Program Signature: %s", randomCode1),
		}),
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	command2 := []string{
		randomExecutableName2, randomName2, randomName3,
	}

	testCase2 := test_cases.CommandWithMultilineResponseTestCase{
		Command: strings.Join(command2, " "),
		MultiLineAssertion: assertions.NewMultiLineAssertion([]string{
			fmt.Sprintf("Program was passed %d args (including program name).", len(command2)),
			fmt.Sprintf("Arg #0 (program name): %s", command2[0]),
			fmt.Sprintf("Arg #1: %s", command2[1]),
			fmt.Sprintf("Arg #2: %s", command2[2]),
			fmt.Sprintf("Program Signature: %s", randomCode2),
		}),
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
