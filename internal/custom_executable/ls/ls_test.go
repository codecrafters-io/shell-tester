package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// Pass the -system flag to use system ls instead of custom implementation
// go test ./... -system
// Tests only pass against BSD implementation of ls, not GNU implementation
// Run on darwin only
var useSystemLs = flag.Bool("system", false, "Use system ls instead of custom implementation")

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

func getLsExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemLs {
		return "ls"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "ls_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "ls_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "ls_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "ls_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS: %s", runtime.GOOS)
	return ""
}

// runLs runs the ls executable with given arguments and returns its output and error if any
func runLs(t *testing.T, args ...string) (string, error) {
	// executable := "ls"
	executable := getLsExecutable(t)

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)
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

func copyLsToDir(t *testing.T, lsFilepath, newDir string) {
	exeData, err := os.ReadFile(lsFilepath)
	if err != nil {
		t.Fatal(err)
	}

	lsFilename := filepath.Base(lsFilepath)
	if err := os.WriteFile(filepath.Join(newDir, lsFilename), exeData, 0755); err != nil {
		t.Fatal(err)
	}
}

func TestLsCurrentDirectory(t *testing.T) {
	if *useSystemLs {
		t.Skip("Skipping test because macOS won't allowing copying coreutils")
	}
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

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
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(oldWd)

	lsExecutable := getLsExecutable(t)
	// Copy ls executable to temp directory
	copyLsToDir(t, lsExecutable, tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Run ls and get output
	output, err := runLs(t)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output
	lsExecutableName := filepath.Base(lsExecutable)
	files := []string{"dir1", "file1.txt", "file2.txt", lsExecutableName}
	// file_name is like ./ls (but in output ./ will not be present)
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
	defer cleanupDirectories([]string{tmpDir})

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

	tmpDir2, err := os.MkdirTemp("", "ls-test2-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir1, tmpDir2})

	// Create test files using the utility function
	createTestFiles(t, tmpDir1, []testFile{
		{name: "file1.txt", content: []byte("test"), mode: 0644},
	})
	createTestFiles(t, tmpDir2, []testFile{
		{name: "file2.txt", content: []byte("test"), mode: 0644},
	})

	// Run ls and get sorted output
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

func TestLsMultipleDirectoriesWithNonExistent(t *testing.T) {
	// Create temporary directories for testing
	tmpDir, err := os.MkdirTemp("", "ls-test1-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	// Create test files using the utility function
	createTestFiles(t, tmpDir, []testFile{
		{name: "file1.txt", content: []byte("test"), mode: 0644},
	})

	// Run ls and get sorted output
	output, _ := runLs(t, tmpDir, "xenon", tmpDir, "non", tmpDir, "bon")

	// Verify output
	expectedOutput := []string{
		"ls: bon: No such file or directory\n",
		"ls: non: No such file or directory\n",
		"ls: xenon: No such file or directory\n",
		tmpDir + ":\nfile1.txt\n\n",
		tmpDir + ":\nfile1.txt\n\n",
		tmpDir + ":\nfile1.txt\n",
	}
	if output != strings.Join(expectedOutput, "") {
		t.Errorf("Expected output %q, got %q", strings.Join(expectedOutput, ""), output)
	}
}

func TestLsWithDashOneFlag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

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
	if *useSystemLs {
		t.Skip("Skipping test because system ls is used")
	}
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ls-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDirectories([]string{tmpDir})

	output, err := runLs(t, "-n", tmpDir)
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "ls: invalid option") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -n") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}

func TestLsWithUnsupportedFlag2(t *testing.T) {
	if *useSystemLs {
		t.Skip("Skipping test because system ls is used")
	}

	// Run ls and get output
	output, err := runLs(t, "-n")
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "ls: invalid option") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -n") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}

func TestLsWithUnsupportedFlag3(t *testing.T) {
	if *useSystemLs {
		t.Skip("Skipping test because system ls is used")
	}

	// Run ls and get output
	output, err := runLs(t, "-l -a")
	if err == nil {
		t.Error("Expected error for unsupported flag, got none")
	}

	if !strings.Contains(output, "ls: invalid option") {
		t.Errorf("Expected internal error notification")
	}
	// Verify error message
	if !strings.Contains(output, "flag provided but not defined: -l -a") {
		t.Errorf("Expected error about undefined flag, got: %q", output)
	}
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
