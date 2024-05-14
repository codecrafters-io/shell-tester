package internal

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testTest(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	logger.Successf("Setup complete")

	b.FeedStdin([]byte("ls"))

	res, err := b.Result()
	if err != nil {
		return err
	}
	logger.Successf(strings.TrimSpace(string(res.Stdout)))
	logger.Infof("Has exited %s", strconv.FormatBool(b.HasExited()))

	b.FeedStdin([]byte("exit"))
	res, err = b.Result()
	if err != nil {
		return err
	}
	logger.Infof("Has exited %s", strconv.FormatBool(b.HasExited()))
	return nil
}
