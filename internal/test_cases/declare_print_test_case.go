package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// DeclarePrintTestCase tests `declare -p VAR` when the variable exists.
// Expected output format: declare -- VAR="value"
type DeclarePrintTestCase struct {
	Variable string
	Value    string
}

func (t DeclarePrintTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := CommandResponseTestCase{
		Command:        fmt.Sprintf("declare -p %s", t.Variable),
		ExpectedOutput: fmt.Sprintf(`declare -- %s="%s"`, t.Variable, t.Value),
		FallbackPatterns: []*regexp.Regexp{
			// zsh uses `typeset VAR=value` format
			regexp.MustCompile(fmt.Sprintf(`^typeset %s=%s$`, regexp.QuoteMeta(t.Variable), regexp.QuoteMeta(t.Value))),
		},
		SuccessMessage: "✓ Received expected response",
	}
	return testCase.Run(asserter, shell, logger)
}
