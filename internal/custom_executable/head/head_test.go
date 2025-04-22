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
	fmt.Printf("=== RUN:  > head %s\n", strings.Join(args, " "))
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
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
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
	output, exitCode, err := runHead(t, "test.txt")
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
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
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
	output, exitCode, err := runHead(t, "-n", "5", "test.txt")
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
	output, exitCode, err = runHead(t, "--lines=3", "test.txt")
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
	output, exitCode, err = runHead(t, "-n", "-2", "test.txt")
	if err == nil {
		t.Fatalf("Expected error, got no error")
	}

	if !strings.Contains(output, "head: illegal line count -- -2") {
		t.Errorf("Expected error to contain 'head: illegal line count -- -2', got: %s", output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestHeadWithByteCount(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "head-test-*")
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

	// Create a test file
	content := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Test with -c flag
	output, exitCode, err := runHead(t, "-c", "10", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that only the first 10 bytes are printed
	expected := "ABCDEFGHIJ"
	if output != expected {
		t.Errorf("Expected output to contain first 10 bytes, got: %q", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Test with --bytes flag
	output, exitCode, err = runHead(t, "--bytes=5", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that only the first 5 bytes are printed
	expected = "ABCDE"
	if output != expected {
		t.Errorf("Expected output to contain first 5 bytes, got: %q", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Test with -c and negative value
	output, exitCode, err = runHead(t, "-c", "-16", "test.txt")
	if err == nil {
		t.Fatalf("Expected error, got no error")
	}

	if !strings.Contains(output, "head: illegal byte count -- -16") {
		t.Errorf("Expected error to contain 'head: illegal byte count -- -16', got: %s", output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestHeadWithMultipleFiles(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "head-test-*")
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

	// Create test files
	testFiles := []testFile{
		{name: "file1.txt", content: []byte("A\nB\nC\nD\nE\n"), mode: 0644},
		{name: "file2.txt", content: []byte("1\n2\n3\n4\n5\n"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Test with multiple files
	output, exitCode, err := runHead(t,
		"file1.txt",
		"file2.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that headers are printed for each file
	file1Path := "file1.txt"
	file2Path := "file2.txt"
	expected := fmt.Sprintf("==> %s <==\nA\nB\nC\nD\nE\n\n==> %s <==\n1\n2\n3\n4\n5\n",
		file1Path, file2Path)

	if output != expected {
		t.Errorf("Expected output with headers, got:\n%s", output)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestHeadWithNonExistentFile(t *testing.T) {
	// Run head on a non-existent file
	output, exitCode, _ := runHead(t, "nonexistent.txt")

	// Check that an error message is printed
	if !strings.Contains(output, "nonexistent.txt") || !strings.Contains(output, "No such file") {
		t.Errorf("Expected error message about non-existent file, got: %s", output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// func TestHeadWithStdin(t *testing.T) {
// 	executable := getHeadExecutable(t)

// 	// Create a pipe to feed input to head
// 	cmd := exec.Command("echo", "-e", "A\nB\nC\nD\nE\nF\nG\nH\nI\nJ\nK\nL")
// 	headCmd := exec.Command(executable, "-n", "5")
// 	headCmd.Stdin, _ = cmd.StdoutPipe()

// 	fmt.Println("=== RUN:  > echo -e \"A\\nB\\nC\\nD\\nE\\nF\\nG\\nH\\nI\\nJ\\nK\\nL\" | head -n 5")

// 	output, err := headCmd.Output()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got: %v", err)
// 	}

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("Failed to start echo command: %v", err)
// 	}

// 	expected := "A\nB\nC\nD\nE\n"
// 	if string(output) != expected {
// 		t.Errorf("Expected first 5 lines from stdin, got: %q", string(output))
// 	}
// }
