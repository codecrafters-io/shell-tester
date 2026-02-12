package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/tester-utils/random"
)

var SMALL_WORDS = []string{"ant", "bee", "cow", "dog", "fox", "owl", "pig", "rat"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

const CUSTOM_CAT_COMMAND = "cat"
const CUSTOM_GREP_COMMAND = "grep"
const CUSTOM_HEAD_COMMAND = "head"
const CUSTOM_LS_COMMAND = "ls"
const CUSTOM_TAIL_COMMAND = "tail"
const CUSTOM_WC_COMMAND = "wc"
const CUSTOM_YES_COMMAND = "yes"

type testCaseContent struct {
	Input            string
	ExpectedOutput   string
	FallbackPatterns []*regexp.Regexp
}

func newTestCaseContent(input string, expectedOutput string, fallbackPatterns []*regexp.Regexp) testCaseContent {
	return testCaseContent{
		Input:            input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: fallbackPatterns,
	}
}

func newTestCaseContents(inputs []string, expectedOutputs []string) []testCaseContent {
	testCases := []testCaseContent{}
	for i, input := range inputs {
		testCases = append(testCases, newTestCaseContent(input, expectedOutputs[i], nil))
	}
	return testCases
}

func newTestCaseContentsWithFallbackPatterns(inputs []string, expectedOutputs []string, fallbackPatterns [][]*regexp.Regexp) []testCaseContent {
	testCases := []testCaseContent{}
	for i, input := range inputs {
		testCases = append(testCases, newTestCaseContent(input, expectedOutputs[i], fallbackPatterns[i]))
	}
	return testCases
}

func getRandomInvalidCommand() string {
	return getRandomInvalidCommands(1)[0]
}

func getRandomInvalidCommands(n int) []string {
	words := random.RandomWords(n)
	invalidCommands := make([]string, n)

	for i := 0; i < n; i++ {
		invalidCommands[i] = "invalid_" + words[i] + "_command"
	}

	return invalidCommands
}

func getRandomString() string {
	// We will use a random numeric string of length = 6
	var result string
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func getRandomName() string {
	names := []string{"Alice", "David", "Emily", "James", "Maria"}
	return names[random.RandomInt(0, len(names))]
}

func logAndQuit(asserter *logged_shell_asserter.LoggedShellAsserter, err error) error {
	asserter.LogRemainingOutput()
	return err
}

func GetRandomCommandSuitableForFile() string {
	return random.RandomElementFromArray([]string{
		"cat",
		"stat",
		"du",
		"wc",
	})
}
