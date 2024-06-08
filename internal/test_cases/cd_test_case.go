package test_cases

import (
	"fmt"
	"os"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type CDAndPWDTestCase struct {
	Directory string // Relative Path possibly
	Response  string // Absolute Path
}

func (t *CDAndPWDTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// First we make sure the directory exists, if not we create it
	command := fmt.Sprintf("cd %s", t.Directory)
	_, err := os.Stat(t.Response)
	if err != nil {
		err = os.Mkdir(t.Response, 0755)
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error creating tmp directory: %v", err)
		}
	}

	// Then we check if prompt is printed
	promptTestCase := NewPromptTestCase("$ ")
	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}
	// And send the cd command, we don't expect any response
	if err := shell.SendCommand(command); err != nil {
		return err
	}

	nextCommand := "pwd"

	// Next we send pwd and check that the directory we cd'ed into is the response
	testCase := SingleLineOutputTestCase{
		Command:                    nextCommand,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`^%s\r\n`, t.Response)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", t.Response+"\n"),
		SuccessMessage:             "Received current working directory response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return nil
}
