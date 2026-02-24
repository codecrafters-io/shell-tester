package shell_executable

import (
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/stretchr/testify/assert"
	"go.chromium.org/luci/common/system/environ"
)

func TestMemoryLimit(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Memory limiting is only supported on Linux")
	}

	path := "test_helpers/memory_hog.sh"
	stageLogger := logger.GetQuietLogger("")
	e := &ShellExecutable{
		executable:         executable.NewExecutable(path),
		stageLogger:        stageLogger,
		programLogger:      logger.GetLogger(stageLogger.IsDebug, "[your-program] "),
		MemoryLimitInBytes: defaultMemoryLimitBytes,
	}
	e.env = environ.New(os.Environ())

	// Set a 50MB memory limit
	e.MemoryLimitInBytes = 50 * 1024 * 1024

	err := e.Start()
	assert.NoError(t, err)

	exited, exitCode := e.WaitForTermination()
	assert.True(t, exited, "Process should have terminated (OOM killed)")
	assert.Equal(t, 137, exitCode, "Expected exit code to be 137 (SIGKILL)")

	err = e.MemoryLimitExceededError()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrMemoryLimitExceeded), "Expected ErrMemoryLimitExceeded, got: %v", err)
	assert.Contains(t, err.Error(), "50 MB", "Error message should contain human-readable memory limit")
}
