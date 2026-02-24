package shell_executable

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryLimit(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Memory limiting is only supported on Linux")
	}

	// Path relative to test package dir (internal/shell_executable)
	e := NewShellExecutableForTest("../test_helpers/oom/memory_hog.sh", nil)
	// Set a 50MB memory limit
	e.MemoryLimitInBytes = 50 * 1024 * 1024

	err := e.Start()
	assert.NoError(t, err)

	exited, _ := e.WaitForTermination()
	assert.True(t, exited, "Process should have terminated (OOM killed)")

	assert.True(t, e.WasOOMKilled(), "Expected process to be OOM killed")

	err = e.MemoryLimitExceededError()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrMemoryLimitExceeded), "Expected ErrMemoryLimitExceeded, got: %v", err)
	assert.Contains(t, err.Error(), "50 MB", "Error message should contain human-readable memory limit")
}
