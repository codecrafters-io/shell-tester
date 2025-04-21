package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Pass the -system flag to use system head instead of custom implementation
var useSystemHead = flag.Bool("system", false, "Use system head instead of custom implementation")

// testFile represents a file to be created for testing
type testFile struct {
	name    string
	content []byte
	mode    os.FileMode
}

// createTestFiles creates test files in the given directory
func createTestFiles(t *testing.T, dir string, files []testFile) {
	t.Helper()
	for _, f := range files {
		path := filepath.Join(dir, f.name)
		if err := os.WriteFile(path, f.content, f.mode); err != nil {
			t.Fatal(err)
		}
	}
}

func getHeadExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemHead {
		return "head"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "head_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "head_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "head_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "head_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runHead runs the head executable with given arguments and returns its output and error if any
func runHead(t *testing.T, args ...string) (string, int, error) {
	executable := getHeadExecutable(t)

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)
	output, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if _, ok := err.(*exec.ExitError); ok {
			exitError = err.(*exec.ExitError)
			exitCode = exitError.ExitCode()
		}
	}
	return string(output), exitCode, err
}

func prettyPrintCommand(args []string) {
	fileParts := strings.Split(args[len(args)-1], "/")
	fileName := fileParts[len(fileParts)-1]

	fmt.Printf("=== RUN:  > head %s\n", strings.Join(args[0:len(args)-1], " ")+" "+fileName)
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func TestHeadDefaultBehavior(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "head-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	// Create a test file with more than 10 lines
	content := strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
		"Line 6",
		"Line 7",
		"Line 8",
		"Line 9",
		"Line 10",
		"Line 11",
		"Line 12",
	}, "\n") + "\n"

	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run head on the test file
	output, exitCode, err := runHead(t, filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that only the first 10 lines are printed
	expected := strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
		"Line 6",
		"Line 7",
		"Line 8",
		"Line 9",
		"Line 10",
	}, "\n") + "\n"

	if output != expected {
		t.Errorf("Expected output to contain first 10 lines, got:\n%s", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestHeadWithLineCount(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "head-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	// Create a test file with more than 5 lines
	content := strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
		"Line 6",
		"Line 7",
	}, "\n") + "\n"

	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Test with -n flag
	output, exitCode, err := runHead(t, "-n", "5", filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that only the first 5 lines are printed
	expected := strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
	}, "\n") + "\n"

	if output != expected {
		t.Errorf("Expected output to contain first 5 lines, got:\n%s", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Test with --lines flag
	output, exitCode, err = runHead(t, "--lines=3", filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that only the first 3 lines are printed
	expected = strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
	}, "\n") + "\n"

	if output != expected {
		t.Errorf("Expected output to contain first 3 lines, got:\n%s", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Test with -n and negative value
	output, exitCode, err = runHead(t, "-n", "-2", filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that all but the last 2 lines are printed
	expected = strings.Join([]string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
	}, "\n") + "\n"

	if output != expected {
		t.Errorf("Expected output to contain all but last 2 lines, got:\n%s", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

}
