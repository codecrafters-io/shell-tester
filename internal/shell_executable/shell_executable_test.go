package shell_executable

// skip running for now for debugging
//
// func TestMemoryLimit(t *testing.T) {
// 	if runtime.GOOS != "linux" {
// 		t.Skip("Memory limiting is only supported on Linux")
// 	}

// 	// Path relative to test package dir (internal/shell_executable)
// 	e := newShellExecutableForTest("test_helpers/memory_hog.sh", nil)
// 	// Set a 50MB memory limit
// 	e.MemoryLimitInBytes = 50 * 1024 * 1024

// 	err := e.Start()
// 	assert.NoError(t, err)

// 	exited, exitCode := e.WaitForTermination()
// 	assert.True(t, exited, "Process should have terminated (OOM killed)")
// 	assert.Equal(t, 137, exitCode, "Expected exit code to be 137 (SIGKILL)")

// 	err = e.MemoryLimitExceededError()
// 	assert.NotNil(t, err)
// 	assert.True(t, errors.Is(err, ErrMemoryLimitExceeded), "Expected ErrMemoryLimitExceeded, got: %v", err)
// 	assert.Contains(t, err.Error(), "50 MB", "Error message should contain human-readable memory limit")
// }

// // NewShellExecutableForTest creates a ShellExecutable that runs the executable at path.
// // Used for testing (e.g. memory limit tests). The logger can be nil; a quiet logger will be used.
// func newShellExecutableForTest(path string, stageLogger *logger.Logger) *ShellExecutable {
// 	if stageLogger == nil {
// 		stageLogger = logger.GetQuietLogger("")
// 	}
// 	b := &ShellExecutable{
// 		executable:         executable.NewExecutable(path),
// 		stageLogger:        stageLogger,
// 		programLogger:      logger.GetLogger(stageLogger.IsDebug, "[your-program] "),
// 		MemoryLimitInBytes: defaultMemoryLimitBytes,
// 	}
// 	b.env = environ.New(os.Environ())
// 	return b
// }
