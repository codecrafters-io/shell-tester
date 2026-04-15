package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// DeclarePrintErrorTestCase tests `declare -p VAR` when the variable does not exist.
// Expected output: declare: VAR: not found
type DeclarePrintErrorTestCase struct {
	Variable string
}

func (t DeclarePrintErrorTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	fallbackPatterns := []*regexp.Regexp{
		regexp.MustCompile(fmt.Sprintf(`^bash: declare: %s: not found$`, regexp.QuoteMeta(t.Variable))),
		regexp.MustCompile(fmt.Sprintf(`^declare: no such variable: %s$`, regexp.QuoteMeta(t.Variable))),
	}

	if !isValidIdentifier(t.Variable) {
		// zsh rejects invalid identifiers passed to -p
		fallbackPatterns = append(fallbackPatterns, regexp.MustCompile(fmt.Sprintf(`^declare: bad argument to -p: %s$`, regexp.QuoteMeta(t.Variable))))
	}

	testCase := CommandResponseTestCase{
		Command:          fmt.Sprintf("declare -p %s", t.Variable),
		ExpectedOutput:   fmt.Sprintf("declare: %s: not found", t.Variable),
		FallbackPatterns: fallbackPatterns,
		SuccessMessage:   "✓ Received expected response",
	}
	return testCase.Run(asserter, shell, logger)
}
