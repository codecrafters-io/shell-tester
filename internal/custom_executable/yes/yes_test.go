package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// Pass the -system flag to use system yes instead of custom implementation
var useSystemYes = flag.Bool("system", false, "Use system yes instead of custom implementation")

func getYesExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemYes {
		return "yes"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "yes_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "yes_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "yes_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "yes_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runYesWithTimeout runs the yes executable with given arguments and returns its output and error if any
// It stops the command after a timeout to prevent the infinite output
func runYesWithTimeout(t *testing.T, timeout time.Duration, args ...string) ([]string, error) {
	executable := getYesExecutable(t)

	t.Helper()
	fmt.Printf("=== RUN:  > yes %s\n", args)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the command with the context
	cmd := exec.CommandContext(ctx, executable, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Read output with a separate timeout to ensure we don't block indefinitely
	lines := []string{}
	scanner := bufio.NewScanner(stdout)

	// Read up to 10 lines or until timeout
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		count := 0
		for scanner.Scan() && count < 40_000 {
			lines = append(lines, scanner.Text())
			count++
		}
	}()

	// Wait for either context timeout or reading to complete
	select {
	case <-ctx.Done():
		// Context timed out, ensure the process is terminated
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	case <-readDone:
		// Reading completed
	}

	// Wait but ignore error as we expect the process might be killed
	cmd.Wait()

	fmt.Printf("=== RUN:  Received %d lines\n", len(lines))
	return lines, nil
}

func TestYesDefault(t *testing.T) {
	lines, err := runYesWithTimeout(t, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(lines) == 0 {
		t.Fatalf("Expected output, got none")
	}

	for _, line := range lines {
		if line != "y" {
			t.Errorf("Expected 'y', got %q", line)
		}
	}
}

func TestYesWithSingleArgument(t *testing.T) {
	lines, err := runYesWithTimeout(t, 100*time.Millisecond, "test")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(lines) == 0 {
		t.Fatalf("Expected output, got none")
	}

	for _, line := range lines {
		if line != "test" {
			t.Errorf("Expected 'test', got %q", line)
		}
	}
}
