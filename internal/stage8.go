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

	randomCodes := []string{
		getRandomString(),
		getRandomString(),
	}

	randomExecutableNames := []string{
		"custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999)),
		"custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999)),
	}

	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: randomExecutableNames[0], CommandMetadata: randomCodes[0]},
		{CommandType: "signature_printer", CommandName: randomExecutableNames[1], CommandMetadata: randomCodes[1]},
	}, true)
	if err != nil {
		return err
	}
	logAvailableExecutables(logger, randomExecutableNames)

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	argCounts := random.RandomInts(1, 4, 2)
	for i, argCount := range argCounts {
		command := []string{randomExecutableNames[i]}
		for range argCount {
			command = append(command, getRandomName())
		}

		expectedLines := []string{
			fmt.Sprintf("Program was passed %d args (including program name).", len(command)),
			fmt.Sprintf("Arg #0 (program name): %s", command[0]),
		}
		for j, arg := range command[1:] {
			expectedLines = append(expectedLines, fmt.Sprintf("Arg #%d: %s", j+1, arg))
		}
		expectedLines = append(expectedLines, fmt.Sprintf("Program Signature: %s", randomCodes[i]))

		testCase := test_cases.CommandWithMultilineResponseTestCase{
			Command:            strings.Join(command, " "),
			MultiLineAssertion: assertions.NewMultiLineAssertion(expectedLines),
			SuccessMessage:     "✓ Received expected response",
		}
		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	return logAndQuit(asserter, nil)
}
