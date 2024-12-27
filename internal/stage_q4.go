package internal

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	randomDir, err := getShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	randomUniqueFileNames := random.RandomInts(1, 100, 3)
	L := random.RandomElementsFromArray(LARGE_WORDS, 6)
	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf(`'f %d'`, randomUniqueFileNames[0])),
		path.Join(randomDir, fmt.Sprintf(`'f  \%d'`, randomUniqueFileNames[1])),
		path.Join(randomDir, fmt.Sprintf(`'f \%d\'`, randomUniqueFileNames[2])),
	}
	fileContents := []string{
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + ".",
		strings.Join(random.RandomWords(2), " ") + "." + "\n",
	}

	inputs := []string{
		fmt.Sprintf(`echo '%s\\\n%s'`, L[0], L[1]),
		fmt.Sprintf(`echo '%s\"%s%s\"%s'`, L[2], L[3], L[4], L[0]),
		fmt.Sprintf(`echo '%s\\n%s'`, L[4], L[1]),
		fmt.Sprintf(`cat "%s" "%s" "%s"`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf(`%s\\\n%s`, L[0], L[1]),
		fmt.Sprintf(`%s\"%s%s\"%s`, L[2], L[3], L[4], L[0]),
		fmt.Sprintf(`%s\\n%s`, L[4], L[1]),
		fileContents[0] + fileContents[1] + strings.TrimRight(fileContents[2], "\n"),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
		testCase := test_cases.CommandResponseTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "✓ Received expected response",
		}

		// For single-quoted strings with line continuation
		if strings.Contains(testCaseContent.Input, `\\\n`) {
			parts := strings.Split(testCaseContent.ExpectedOutput, `\\\n`)
			if len(parts) == 2 {
				firstPart := parts[0]
				secondPart := parts[1]

				// Add fallback patterns for both bash and ash output formats
				testCase.FallbackPatterns = []*regexp.Regexp{
					// Pattern for bash-style single line output
					regexp.MustCompile(`^` + regexp.QuoteMeta(testCaseContent.ExpectedOutput) + `$`),
					// Pattern for ash's line-split format
					regexp.MustCompile(`^` + regexp.QuoteMeta(firstPart) + `\\$`),
					regexp.MustCompile(`^` + regexp.QuoteMeta(secondPart) + `$`),
				}
			}
		}

		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	testCase := test_cases.CommandResponseTestCase{
		Command:          testCaseContents[3].Input,
		ExpectedOutput:   testCaseContents[3].ExpectedOutput,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
