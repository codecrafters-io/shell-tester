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
	randomCode := getRandomString()
	randomName := getRandomName()
	randomExecutableName := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: randomExecutableName, CommandMetadata: randomCode},
	}, true)
	if err != nil {
		return err
	}
	logAvailableExecutables(logger, []string{randomExecutableName})
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	command := []string{
		randomExecutableName, randomName,
	}

	testCase := test_cases.CommandWithMultilineResponseTestCase{
		Command: strings.Join(command, " "),
		MultiLineAssertion: assertions.NewMultiLineAssertion([]string{
			fmt.Sprintf("Program was passed %d args (including program name).", len(command)),
			fmt.Sprintf("Arg #0 (program name): %s", command[0]),
			fmt.Sprintf("Arg #1: %s", command[1]),
			fmt.Sprintf("Program Signature: %s", randomCode),
		}),
		SuccessMessage: "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
