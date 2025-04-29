package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Pass the -system flag to use system grep instead of custom implementation
// go test ./... -system
// Tests only pass against BSD implementation of grep, not GNU implementation
// Run on darwin only
var useSystemGrep = flag.Bool("system", false, "Use system grep instead of custom implementation")

func getGrepExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemGrep {
		return "grep"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "grep_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "grep_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "grep_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "grep_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runGrep runs the grep executable with given arguments and returns its output and error if any
func runGrep(t *testing.T, stdinContent string, args ...string) (string, string, int, error) {
	executable := getGrepExecutable(t)

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)

	if stdinContent != "" {
		cmd.Stdin = strings.NewReader(stdinContent)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}

	return stdout.String(), stderr.String(), exitCode, err // Return err as well
}

func prettyPrintCommand(args []string) {
	// Basic pretty printing, similar to cat/wc tests
	displayArgs := make([]string, len(args))
	copy(displayArgs, args)
	// Potentially shorten paths or quote arguments if needed here

	out := fmt.Sprintf("=== RUN:  > grep %s", strings.Join(displayArgs, " "))
	fmt.Println(out)
}

// TestGrepStdin tests grep functionality reading from standard input.
func TestGrepStdin(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		input        string
		expectedOut  string
		expectedErr  string
		expectedExit int
	}{
		{
			name:         "Simple match",
			pattern:      "hello",
			input:        "hello world\nthis is a test\nhello again",
			expectedOut:  "hello world\nhello again\n",
			expectedErr:  "",
			expectedExit: 0,
		},
		{
			name:         "No match",
			pattern:      "goodbye",
			input:        "hello world\nthis is a test",
			expectedOut:  "",
			expectedErr:  "",
			expectedExit: 1,
		},
		{
			name:         "Empty input",
			pattern:      "test",
			input:        "",
			expectedOut:  "",
			expectedErr:  "",
			expectedExit: 1,
		},
		{
			name:         "Match on empty string pattern",
			pattern:      "", // Our simple strings.Contains matches everything
			input:        "line1\nline2",
			expectedOut:  "line1\nline2\n",
			expectedErr:  "",
			expectedExit: 0,
			// Note: Real grep might error or behave differently with empty pattern
		},
		{
			name:         "Usage error - no pattern",
			pattern:      "", // Args will be empty
			input:        "test",
			expectedOut:  "",
			expectedErr:  "Usage: grep PATTERN\n",
			expectedExit: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []string

			if *useSystemGrep && tt.name == "Usage error - no pattern" {
				return // Skip this test for System grep
			}
			if tt.name != "Usage error - no pattern" {
				args = append(args, tt.pattern)
			} // else args remains empty

			stdout, stderr, exitCode, runErr := runGrep(t, tt.input, args...)

			// Check exit code first, as stderr/stdout might be irrelevant if exit code is wrong
			if exitCode != tt.expectedExit {
				// Include stdout/stderr in the error message for better debugging
				t.Errorf("Expected exit code %d, got %d. Stderr: %q, Stdout: %q, RunErr: %v",
					tt.expectedExit, exitCode, stderr, stdout, runErr)
			}

			// Now check stdout and stderr
			if stdout != tt.expectedOut {
				t.Errorf("Expected stdout %q, got %q", tt.expectedOut, stdout)
			}

			if stderr != tt.expectedErr {
				t.Errorf("Expected stderr %q, got %q", tt.expectedErr, stderr)
			}

			// Handle expected usage error specifically
			if tt.expectedExit == 2 && runErr == nil {
				t.Errorf("Expected a non-nil error for usage error case, but got nil")
			} else if tt.expectedExit != 2 && runErr != nil {
				// If we didn't expect an error exit, but got one, check if it was *exec.ExitError
				if _, ok := runErr.(*exec.ExitError); !ok {
					// It was some other error during execution
					t.Errorf("Command execution failed unexpectedly: %v", runErr)
				}
				// If it *was* an ExitError, we already checked the exit code above.
			}
		})
	}
}
