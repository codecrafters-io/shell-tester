package internal

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testP3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "grep", CommandName: CUSTOM_GREP_COMMAND, CommandMetadata: ""},
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
		{CommandType: "ls", CommandName: CUSTOM_LS_COMMAND, CommandMetadata: ""},
		{CommandType: "wc", CommandName: CUSTOM_WC_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Test-1
	randomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	filePath := path.Join(randomDir, fmt.Sprintf("file-%d", random.RandomInt(1, 100)))
	randomWords := random.RandomWords(5)
	firstThreeLines := fmt.Sprintf("%s\n%s\n%s\n", randomWords[0], randomWords[1], randomWords[2])
	fileContent := firstThreeLines + fmt.Sprintf("%s\n%s\n", randomWords[3], randomWords[4])
	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	lines := strings.Count(firstThreeLines, "\n")
	words := strings.Count(strings.ReplaceAll(firstThreeLines, "\n", " "), " ")
	bytes := len(firstThreeLines)

	input := fmt.Sprintf(`cat %s | head -n 3 | wc`, filePath)
	expectedOutput := fmt.Sprintf("%8d%8d%8d", lines, words, bytes)

	singleLineTestCase := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test-2
	newRandomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	randomUniqueFileNames := random.RandomInts(1, 100, 6)
	filePaths := []string{
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[0])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[1])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[2])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[3])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[4])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[5])),
	}
	fileContents := random.RandomWords(6)
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	sort.Slice(randomUniqueFileNames, func(i, j int) bool {
		a, b := strconv.Itoa(randomUniqueFileNames[i]), strconv.Itoa(randomUniqueFileNames[j])
		return a < b
	})
	availableEntries := randomUniqueFileNames[1:4]

	input = fmt.Sprintf(`ls %s | tail -n 5 | head -n 3 | grep "f-%d"`, newRandomDir, availableEntries[2])
	expectedOutput = fmt.Sprintf("f-%d", availableEntries[2])
	expectedRegexPattern := fmt.Sprintf("^f-%d$", availableEntries[2])

	singleLineTestCase2 := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(expectedRegexPattern)},
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase2.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
