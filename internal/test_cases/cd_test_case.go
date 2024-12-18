package test_cases

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type CDAndPWDTestCase struct {
	Directory string // Relative Path possibly
	Response  string // Absolute Path
}

func (t *CDAndPWDTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// First we make sure the directory exists, if not we create it
	command := fmt.Sprintf("cd %s", t.Directory)
	_, err := os.Stat(t.Response)
	if err != nil {
		err = os.Mkdir(t.Response, 0755)
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error creating tmp directory: %v", err)
		}
	}

	// And send the cd command, we don't expect any response
	tc := CommandReflectionTestCase{
		Command:             command,
		SkipPromptAssertion: false,
	}
	if err := tc.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	nextCommand := "pwd"

	// Next we send pwd and check that the directory we cd'ed into is the response
	testCase := CommandResponseTestCase{
		Command:          nextCommand,
		ExpectedOutput:   t.Response,
		FallbackPatterns: nil,
		SuccessMessage:   "Received current working directory response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}
