package test_cases

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// HistoryPersistenceTestCase checks that after loading history from a file, the shell's history output matches the file's contents.
type HistoryPersistenceTestCase struct {
	// PreviousCommands are the commands that were executed before loading the history file
	PreviousCommands []CommandResponseTestCase

	// FilePath is the path to the file containing history commands
	FilePath string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// Was history command executed before loading the history file
	WasHistoryCommandExecuted bool
}

func (t HistoryPersistenceTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	// Read the file contents
	file, err := os.Open(t.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open history file: %v", err)
	}
	defer file.Close()

	var fileCommands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			fileCommands = append(fileCommands, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	// Send the 'history' command to the shell
	if err := shell.SendCommand("history"); err != nil {
		return fmt.Errorf("failed to send history command: %v", err)
	}

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ history",
	})

	// Check that each previous command appears in the history output
	for i, cmd := range t.PreviousCommands {
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   fmt.Sprintf("    %d  %s", i+1, cmd.Command),
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^\s*\d+\s+%s$`, regexp.QuoteMeta(cmd.Command)))},
		})
	}

	// Add assertion for the first history command if it was executed
	if t.WasHistoryCommandExecuted {
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   fmt.Sprintf("    %d  history", len(t.PreviousCommands)+1),
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(`^\s*\d+\s+history$`)},
		})
	}

	// Check that each file command appears in the history output
	historyOffset := 1
	if t.WasHistoryCommandExecuted {
		historyOffset = 2
	}

	// Add assertion for the history -r command first
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   fmt.Sprintf("    %d  history -r %s", len(t.PreviousCommands)+historyOffset, t.FilePath),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^\s*\d+\s+history -r %s$`, regexp.QuoteMeta(t.FilePath)))},
	})

	// Then check file commands
	for i, cmd := range fileCommands {
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   fmt.Sprintf("    %d  %s", len(t.PreviousCommands)+historyOffset+1+i, cmd),
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^\s*\d+\s+%s$`, regexp.QuoteMeta(cmd)))},
		})
	}

	// Add assertion for the final history command
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   fmt.Sprintf("    %d  history", len(t.PreviousCommands)+len(fileCommands)+historyOffset+2),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(`^\s*\d+\s+history$`)},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
