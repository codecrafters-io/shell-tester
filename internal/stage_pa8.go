package internal

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// gitMultiSubcommandChoice drives top-level `git <subcommand>` completion only (no nested words).
// Completion strings are chosen so bash does not extend a longest-common-prefix past what the
// user already typed: either they share no common prefix (empty partial), or their LCP equals partial.
type gitMultiSubcommandChoice struct {
	partial     string
	completions []string
}

func testPA8(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	command := "git"

	// The partial script should not accidentally trigger LCP for this stage
	// So we choose completions with no common remaining prefix
	choices := []gitMultiSubcommandChoice{
		{partial: "che", completions: []string{"checkout", "cherry-pick"}},
		{partial: "re", completions: []string{"reset", "remote", "rebase"}},
		{partial: "pu", completions: []string{"push", "pull"}},
		{partial: "sta", completions: []string{"stash", "status"}},
	}

	choice := choices[random.RandomInt(0, len(choices))]

	sortedCompletions := append([]string(nil), choice.completions...)
	sort.Strings(sortedCompletions)

	compLineEnvVar := fmt.Sprintf("%s %s", command, choice.partial)
	completionsLine := strings.Join(sortedCompletions, "  ")

	completerPath := path.Join(randomDir, "multiCandidateEnvCompleter")

	secretValue := getRandomString()
	if err := (&custom_executable.CompleterExecutableSpecification{
		Path:        completerPath,
		SecretValue: secretValue,
		CompleterConfiguration: completer_configuration.CompleterConfiguration{
			OutputLines: sortedCompletions,
			ExpectedArguments: &completer_configuration.CompleterConfigurationExpectedArguments{
				Argv1: command,
				Argv2: choice.partial,
				Argv3: command,
			},
			ExpectedEnvVars: &completer_configuration.CompleterConfigurationExpectedEnvVars{
				CompLine:  compLineEnvVar,
				CompPoint: strconv.Itoa(len(compLineEnvVar)),
			},
		},
	}).Create(); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	registerCmd := fmt.Sprintf("complete -C %s %s", completerPath, command)
	registerTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ Registered command-based completion",
	}
	if err := registerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	escapedCompletionChoices := make([]string, len(sortedCompletions))
	for i, name := range sortedCompletions {
		escapedCompletionChoices[i] = regexp.QuoteMeta(name)
	}

	multiCase := test_cases.MultipleCompletionsTestCase{
		RawInput:                      compLineEnvVar,
		TabCount:                      2,
		CheckForBell:                  true,
		ExpectedCompletionOptionsLine: completionsLine,
		ExpectedCompletionOptionsLineFallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile("^" + strings.Join(escapedCompletionChoices, `\s+`) + "$"),
		},
		SkipPromptAssertion: true,
		SuccessMessage:      "✓ Multiple programmable completion options listed",
	}

	if err := multiCase.Run(asserter, shell, logger); err != nil {
		completer_configuration.LogCompleterErrors(logger, secretValue)
		return err
	}

	return logAndQuit(asserter, nil)
}
