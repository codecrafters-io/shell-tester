package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// DeclareAssignmentTestCase tests `VAR=value` style assignment via declare.
// If the variable name is valid, no output is expected.
// If invalid, the shell must print an error.
type DeclareAssignmentTestCase struct {
	Variable string
	Value    string
}

func (t DeclareAssignmentTestCase) IsValidAssignment() bool {
	return isValidIdentifier(t.Variable)
}

func (t DeclareAssignmentTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	assignment := fmt.Sprintf("%s=%s", t.Variable, t.Value)

	if t.IsValidAssignment() {
		testCase := CommandWithNoResponseTestCase{
			Command:        fmt.Sprintf("declare %s", assignment),
			SuccessMessage: "✓ Received expected response",
		}
		return testCase.Run(asserter, shell, logger, false)
	}

	expectedOutput := fmt.Sprintf("declare: `%s': not a valid identifier", assignment)
	testCase := CommandResponseTestCase{
		Command:        fmt.Sprintf("declare %s", assignment),
		ExpectedOutput: expectedOutput,
		FallbackPatterns: []*regexp.Regexp{
			// bash prefixes with "bash: "
			regexp.MustCompile(fmt.Sprintf("^bash: declare: `%s': not a valid identifier$", regexp.QuoteMeta(assignment))),
			// zsh uses "not an identifier" phrasing
			regexp.MustCompile(fmt.Sprintf(`^declare: not an identifier: %s$`, regexp.QuoteMeta(t.Variable))),
		},
		SuccessMessage: "✓ Received expected response",
	}
	return testCase.Run(asserter, shell, logger)
}
