package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// testFile represents a file or directory to be created for testing
type testFile struct {
	name    string
	isDir   bool
	content []byte
	mode    os.FileMode
}

// createTestFiles creates test files in the given directory and returns cleanup function
func createTestFiles(t *testing.T, dir string, files []testFile) {
	t.Helper()
	for _, f := range files {
		path := filepath.Join(dir, f.name)
		if f.isDir {
			if err := os.Mkdir(path, f.mode); err != nil {
				t.Fatal(err)
			}
		} else {
			if err := os.WriteFile(path, f.content, f.mode); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// runLs runs the ls executable with given arguments and returns its output and error if any
func runLs(t *testing.T, args ...string) (string, error) {
	var executable string
	switch runtime.GOOS {
	case "darwin":
		executable = "./ls_darwin_arm64"
	case "linux":
		executable = "./ls_linux_amd64"
	default:
		t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)
	// fmt.Println(cmd.String())
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func prettyPrintCommand(args []string) {
	copiedArgs := make([]string, len(args))
	copy(copiedArgs, args)
	for i, arg := range copiedArgs {
		if !strings.HasPrefix(arg, "-") {
			copiedArgs[i] = strings.Split(arg, "/")[len(strings.Split(arg, "/"))-1]
		}
	}

	executable := "ls"
	out := fmt.Sprintf("=== RUN:  > %s %s", executable, strings.Join(copiedArgs, " "))
	fmt.Println(out)
}

func copyLsToDir(t *testing.T, ls_file, newDir string) {
	exeData, err := os.ReadFile(ls_file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(newDir, ls_file), exeData, 0755); err != nil {
		t.Fatal(err)
	}
}

func TestLsCurrentDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files using the utility function
	testFiles := []testFile{
		{name: "file1.txt", content: []byte("test"), mode: 0644},
		{name: "file2.txt", content: []byte("test"), mode: 0644},
		{name: "dir1", isDir: true, mode: 0755},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Change to the test directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	// Copy ls executable to temp directory
	copyLsToDir(t, "ls_darwin_arm64", tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Run ls and get output
	output, err := runLs(t)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output
	files := []string{"dir1", "file1.txt", "file2.txt", "ls_darwin_arm64"}
	sort.Strings(files)
	expected := strings.Join(files, "\n") + "\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestLsSpecificDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files using the utility function
	testFiles := []testFile{
		{name: "a.txt", content: []byte("test"), mode: 0644},
		{name: "b.txt", content: []byte("test"), mode: 0644},
		{name: "c.txt", content: []byte("test"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run ls and get output
	output, err := runLs(t, tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output
	files := []string{"a.txt", "b.txt", "c.txt"}
	sort.Strings(files)
	expected := strings.Join(files, "\n") + "\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestLsMultipleDirectories(t *testing.T) {
	// Create temporary directories for testing
	tmpDir1, err := os.MkdirTemp("", "ls-test1-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir1)

	tmpDir2, err := os.MkdirTemp("", "ls-test2-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir2)

	// Create test files using the utility function
	createTestFiles(t, tmpDir1, []testFile{
		{name: "file1.txt", content: []byte("test"), mode: 0644},
	})
	createTestFiles(t, tmpDir2, []testFile{
		{name: "file2.txt", content: []byte("test"), mode: 0644},
	})

	// Run ls and get output
	// Output should also be sorted
	output, err := runLs(t, tmpDir2, tmpDir1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output
	expected := tmpDir1 + ":\nfile1.txt\n\n" + tmpDir2 + ":\nfile2.txt\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestLsNonExistentDirectory(t *testing.T) {
	// Run ls and get output
	output, _ := runLs(t, "nonexistent")

	// Verify output contains error message
	expectedError := "ls: nonexistent: No such file or directory\n"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Expected error message containing %q, got %q", expectedError, output)
	}
}

func TestLsNonExistentDirectory2(t *testing.T) {
	// Run ls and get output
	output, _ := runLs(t, "-1", "nonexistent")

	// Verify output contains error message
	expectedError := "ls: nonexistent: No such file or directory\n"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Expected error message containing %q, got %q", expectedError, output)
	}
}

func TestLsNonExistentDirectory3(t *testing.T) {
	// Run ls and get output
	output, _ := runLs(t, "nonexistent", "nonexistent")

	// Verify output contains error message
	expectedError := "ls: nonexistent: No such file or directory\nls: nonexistent: No such file or directory\n"

	if !strings.Contains(output, expectedError) {
		t.Errorf("Expected error message containing %q, got %q", expectedError, output)
	}
}

func TestLsWithDashOneFlag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files using the utility function
	testFiles := []testFile{
		{name: "file1.txt", content: []byte("test"), mode: 0644},
		{name: "file2.txt", content: []byte("test"), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run ls and get output
	output, err := runLs(t, "-1", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output (should be one file per line)
	files := []string{"file1.txt", "file2.txt"}
	sort.Strings(files)
	expected := strings.Join(files, "\n") + "\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestLsWithUnsupportedFlag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	output, err := runLs(t, "-n", tmpDir)
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "CodeCrafters Internal Error") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -n") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}

func TestLsWithUnsupportedFlag2(t *testing.T) {
	// Run ls and get output
	output, err := runLs(t, "-n")
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "CodeCrafters Internal Error") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -n") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}

func TestLsWithUnsupportedFlag3(t *testing.T) {
	// Run ls and get output
	output, err := runLs(t, "-l -a")
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "CodeCrafters Internal Error") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -l -a") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}
