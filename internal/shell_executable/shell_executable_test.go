package shell_executable

import (
	"errors"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/stretchr/testify/assert"
	"go.chromium.org/luci/common/system/environ"
)

func TestMemoryLimit(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Memory limiting is only supported on Linux")
	}

	shell := newMemoryHogExecutable()
	shell.MemoryLimitInBytes = 50 * 1024 * 1024 // 50MB

	err := shell.Start()
	assert.NoError(t, err)

	conditionFn := func() bool {
		return false
	}

	err = shell.ReadUntilConditionOrTimeout(conditionFn, 2000*time.Millisecond)
	assert.True(t, errors.Is(err, ErrProgramExited), "Expected ErrProgramExited, got: %v", err)

	assert.True(t, shell.WasOOMKilled(), "expected OOM killed")
	err = shell.MemoryLimitExceededError()
	assert.True(t, errors.Is(err, ErrMemoryLimitExceeded), "Expected ErrMemoryLimitExceeded, got: %v", err)
	assert.Contains(t, err.Error(), "50 MB", "Error message should contain human-readable memory limit")
}

func newMemoryHogExecutable() *ShellExecutable {
	path := "test_helpers/memory_hog.sh"
	stageLogger := logger.GetQuietLogger("")

	shell := &ShellExecutable{
		executable:         executable.NewExecutable(path),
		stageLogger:        stageLogger,
		programLogger:      logger.GetLogger(stageLogger.IsDebug, "[your-program] "),
		MemoryLimitInBytes: defaultMemoryLimitBytes,
	}
	shell.env = environ.New(os.Environ())
	return shell
}
