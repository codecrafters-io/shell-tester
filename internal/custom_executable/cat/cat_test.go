package main

import (
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

// Pass the -system flag to use system cat instead of custom implementation
// go test ./... -system
// Tests only pass against BSD implementation of cat, not GNU implementation
// Run on darwin only
var useSystemCat = flag.Bool("system", false, "Use system cat instead of custom implementation")

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

func getCatExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemCat {
		return "cat"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "cat_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "cat_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "cat_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "cat_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runCat runs the cat executable with given arguments and returns its output and error if any
func runCat(t *testing.T, args ...string) (string, int, error) {
	executable := getCatExecutable(t)

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)
	output, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}
	return string(output), exitCode, err
}

func prettyPrintCommand(args []string) {
	copiedArgs := make([]string, len(args))
	copy(copiedArgs, args)
	for i, arg := range copiedArgs {
		if !strings.HasPrefix(arg, "-") {
			copiedArgs[i] = strings.Split(arg, "/")[len(strings.Split(arg, "/"))-1]
		}
	}

	out := fmt.Sprintf("=== RUN:  > cat %s", strings.Join(copiedArgs, " "))
	fmt.Println(out)
}

func TestCatSingleFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "cat-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	content := "Hello, World!\n"
	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run cat and get output
	output, exitCode, err := runCat(t, filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output != content {
		t.Errorf("Expected output %q, got %q", content, output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestCatMultipleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cat-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	testFiles := []testFile{
		{name: "file1.txt", content: []byte("Content 1\n"), mode: 0644},
		{name: "file2.txt", content: []byte("Content 2\n"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	output, exitCode, err := runCat(t,
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "Content 1\nContent 2\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestCatNonExistentFile(t *testing.T) {
	output, exitCode, _ := runCat(t, "nonexistent.txt")
	expectedError := "cat: nonexistent.txt: No such file or directory\n"
	if output != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestCatStdin(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cat-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	executable := getCatExecutable(t)
	cmd := exec.Command(executable)
	cmd.Stdin = strings.NewReader("Hello from stdin\n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "Hello from stdin\n"
	if string(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, string(output))
	}
}

func TestCatMixedExistingAndNonExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cat-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	testFiles := []testFile{
		{name: "exists.txt", content: []byte("I exist\n"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	output, exitCode, _ := runCat(t,
		filepath.Join(tmpDir, "exists.txt"),
		"nonexistent.txt",
		filepath.Join(tmpDir, "exists.txt"),
		"nonexistent.txt",
		filepath.Join(tmpDir, "exists.txt"),
		"nonexistent.txt",
	)

	expectedContent := "I exist\n"
	expectedError := "cat: nonexistent.txt: No such file or directory\n"
	expected := expectedContent + expectedError + expectedContent + expectedError + expectedContent + expectedError

	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
