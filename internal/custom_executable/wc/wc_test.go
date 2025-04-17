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

// Pass the -system flag to use system wc instead of custom implementation
// go test ./... -system
// Tests only pass against BSD implementation of wc, not GNU implementation
// Run on darwin only
var useSystemWc = flag.Bool("system", false, "Use system wc instead of custom implementation")

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

func getWcExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemWc {
		return "wc"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "wc_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "wc_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "wc_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "wc_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runWc runs the wc executable with given arguments and returns its output and error if any
func runWc(t *testing.T, args ...string) (string, int, error) {
	executable := getWcExecutable(t)

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
			// copiedArgs[i] = strings.Split(arg, "/")[len(strings.Split(arg, "/"))-1]
			copiedArgs[i] = arg
		}
	}

	out := fmt.Sprintf("=== RUN:  > wc %s", strings.Join(copiedArgs, " "))
	fmt.Println(out)
}

func TestWcSingleFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
	defer cleanupDirectories([]string{tmpDir})

	content := "Hello, World\nis a test file.\n"
	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run wc and get output
	output, exitCode, err := runWc(t, "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "       2       6      29 test.txt\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestWcMultipleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	testFiles := []testFile{
		{name: "file1.txt", content: []byte("First file\n"), mode: 0644},
		{name: "file2.txt", content: []byte("Second file\n"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
	defer cleanupDirectories([]string{tmpDir})

	output, exitCode, err := runWc(t,
		"file1.txt",
		"file2.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "       1       2      11 file1.txt\n" +
		"       1       2      12 file2.txt\n" +
		"       2       4      23 total\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestWcWithFlags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	content := "Hello, World\nis a test file.\n"
	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
	defer cleanupDirectories([]string{tmpDir})

	// Test -l flag
	output, _, err := runWc(t, "-l", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	expected := fmt.Sprintf("       2 %s\n", "test.txt")
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	// Test -w flag
	output, _, err = runWc(t, "-w", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	expected = fmt.Sprintf("       6 %s\n", "test.txt")
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	// Test -c flag
	output, _, err = runWc(t, "-c", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	expected = fmt.Sprintf("      29 %s\n", "test.txt")
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	// Test multiple flags
	output, _, err = runWc(t, "-lw", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	expected = fmt.Sprintf("       2       6 %s\n", "test.txt")
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestWcStdin(t *testing.T) {
	executable := getWcExecutable(t)
	cmd := exec.Command(executable)
	cmd.Stdin = strings.NewReader("Hello from stdin\n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "       1       3      17\n"
	if string(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, string(output))
	}
}

func TestWcNonExistentFile(t *testing.T) {
	output, exitCode, _ := runWc(t, "nonexistent.txt")
	expectedError := "wc: nonexistent.txt: open: No such file or directory\n"
	if output != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestWcMixedExistingAndNonExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	testFiles := []testFile{
		{name: "exists.txt", content: []byte("I exist\n"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
	defer cleanupDirectories([]string{tmpDir})

	output, exitCode, _ := runWc(t,
		"exists.txt",
		"nonexistent.txt",
		"exists.txt",
	)

	expectedContent := "       1       2       8 exists.txt\n"
	expectedError := "wc: nonexistent.txt: open: No such file or directory\n"
	expected := expectedError + expectedContent + expectedContent + "       2       4      16 total\n"

	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestWcWithUnsupportedFlag(t *testing.T) {
	if *useSystemWc {
		t.Skip("Skipping test because system wc is used")
	}

	output, exitCode, err := runWc(t, "-n")
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}
	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, got: %d", exitCode)
	}

	if !strings.Contains(output, "wc: invalid option") {
		t.Errorf("Expected error about invalid option, got: %q", output)
	}
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
