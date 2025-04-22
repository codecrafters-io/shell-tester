package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time" // Needed for -f tests eventually
)

// Pass the -system flag to use system tail instead of custom implementation
var useSystemTail = flag.Bool("system", false, "Use system tail instead of custom implementation")

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

func getTailExecutable(t *testing.T) string {
	testerDir := filepath.Join(os.Getenv("TESTER_DIR"), "built_executables")
	if *useSystemTail {
		return "tail"
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return filepath.Join(testerDir, "tail_darwin_arm64")
		case "amd64":
			return filepath.Join(testerDir, "tail_darwin_amd64")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return filepath.Join(testerDir, "tail_linux_amd64")
		case "arm64":
			return filepath.Join(testerDir, "tail_linux_arm64")
		}
	}
	t.Fatalf("Unsupported OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	return ""
}

// runTail runs the tail executable with given arguments and returns its output and error if any
func runTail(t *testing.T, args ...string) (string, int, error) {
	executable := getTailExecutable(t)

	t.Helper()
	prettyPrintCommand(args)
	cmd := exec.Command(executable, args...)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
	exitCode := 0
	var exitError *exec.ExitError
	if err != nil {
		if ok := errors.As(err, &exitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// If it's not an ExitError, it might be something else (e.g., command not found)
			// We still want to report the error but maybe exitCode 1 is appropriate?
			exitCode = 1 // Default to 1 for non-exit errors? Or specific code?
		}
	}
	// Normalize EOL characters for comparison
	outputStr := strings.ReplaceAll(string(output), "\r\n", "\n")
	return outputStr, exitCode, err
}

// runTailWithInput runs tail with specific stdin content
func runTailWithInput(t *testing.T, stdinContent string, args ...string) (string, int, error) {
	executable := getTailExecutable(t)

	t.Helper()
	prettyPrintCommandWithPipe(stdinContent, args)
	cmd := exec.Command(executable, args...)
	cmd.Stdin = strings.NewReader(stdinContent)
	output, err := cmd.CombinedOutput()
	exitCode := 0
	var exitError *exec.ExitError
	if err != nil {
		if ok := errors.As(err, &exitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}
	// Normalize EOL characters for comparison
	outputStr := strings.ReplaceAll(string(output), "\r\n", "\n")
	return outputStr, exitCode, err
}

func prettyPrintCommand(args []string) {
	fmt.Printf("=== RUN:  > tail %s\n", strings.Join(args, " "))
}

func prettyPrintCommandWithPipe(stdinContent string, args []string) {
	// Truncate long stdin for printing
	displayStdin := stdinContent
	if len(displayStdin) > 50 {
		displayStdin = displayStdin[:47] + "..."
	}
	displayStdin = strings.ReplaceAll(displayStdin, "\n", "\\n") // Make newlines visible
	fmt.Printf("=== RUN:  > echo -n %q | tail %s\n", displayStdin, strings.Join(args, " "))
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func setupTestDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "tail-test-*")
	if err != nil {
		t.Fatal(err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	cleanup := func() {
		os.Chdir(currentDir)
		cleanupDirectories([]string{tmpDir})
	}
	return tmpDir, cleanup
}

// --- Test Cases ---

func TestTailDefaultBehavior(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create a test file with more than 10 lines
	lines := make([]string, 15)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %d", i+1)
	}
	content := strings.Join(lines, "\n") + "\n"

	testFiles := []testFile{
		{name: "test.txt", content: []byte(content), mode: 0644},
	}
	createTestFiles(t, tmpDir, testFiles)

	// Run tail on the test file
	output, exitCode, err := runTail(t, "test.txt")
	if err != nil && exitCode != 0 { // Allow err if exitCode is 0
		t.Fatalf("Expected no error with exit code 0, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}

	// Check that only the last 10 lines are printed (Line 6 to Line 15)
	expectedLines := lines[5:] // Index 5 is Line 6
	expected := strings.Join(expectedLines, "\n") + "\n"

	if output != expected {
		t.Errorf("Expected output to contain last 10 lines (6-15), got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestTailWithLineCount(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %d", i+1)
	}
	content := strings.Join(lines, "\n") + "\n"
	testFiles := []testFile{{name: "test.txt", content: []byte(content), mode: 0644}}
	createTestFiles(t, tmpDir, testFiles)

	// Test with -n N (last N lines)
	output, exitCode, err := runTail(t, "-n", "3", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-n 3: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected := "Line 8\nLine 9\nLine 10\n"
	if output != expected {
		t.Errorf("-n 3: Expected last 3 lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-n 3: Expected exit code 0, got %d", exitCode)
	}

	// Test with --lines=N
	output, exitCode, err = runTail(t, "--lines=5", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("--lines=5: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 6\nLine 7\nLine 8\nLine 9\nLine 10\n"
	if output != expected {
		t.Errorf("--lines=5: Expected last 5 lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("--lines=5: Expected exit code 0, got %d", exitCode)
	}

	// Test with -n +N (start from line N)
	output, exitCode, err = runTail(t, "-n", "+7", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-n +7: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 7\nLine 8\nLine 9\nLine 10\n"
	if output != expected {
		t.Errorf("-n +7: Expected lines from 7 onwards, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-n +7: Expected exit code 0, got %d", exitCode)
	}

	// Test with -n +N where N > line count (should print nothing)
	output, exitCode, err = runTail(t, "-n", "+11", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-n +11: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = ""
	if output != expected {
		t.Errorf("-n +11: Expected empty output, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-n +11: Expected exit code 0, got %d", exitCode)
	}

	// Test combined number flag: -n5
	output, exitCode, err = runTail(t, "-n5", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-n5: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 6\nLine 7\nLine 8\nLine 9\nLine 10\n"
	if output != expected {
		t.Errorf("-n5: Expected last 5 lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-n5: Expected exit code 0, got %d", exitCode)
	}

	// Test invalid count
	output, exitCode, err = runTail(t, "-n", "abc", "test.txt")
	if err == nil {
		t.Fatalf("Expected error for invalid count, got none. Output:\n%s", output)
	}
	if !strings.Contains(output, "illegal offset") {
		t.Errorf("Expected error message 'illegal offset', got: %s", output)
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 for invalid count, got %d", exitCode)
	}
}

// Test for -c (bytes) - requires implementation in tail.go first
func TestTailWithByteCount(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	content := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // 26 bytes
	testFiles := []testFile{{name: "test.txt", content: []byte(content), mode: 0644}}
	createTestFiles(t, tmpDir, testFiles)

	// Test with -c N (last N bytes)
	output, exitCode, err := runTail(t, "-c", "10", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-c 10: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected := "QRSTUVWXYZ" // Last 10 bytes
	if output != expected {
		t.Errorf("-c 10: Expected last 10 bytes %q, got: %q", expected, output)
	}
	if exitCode != 0 {
		t.Errorf("-c 10: Expected exit code 0, got %d", exitCode)
	}

	// Test with --bytes=N
	output, exitCode, err = runTail(t, "--bytes=5", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("--bytes=5: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "VWXYZ" // Last 5 bytes
	if output != expected {
		t.Errorf("--bytes=5: Expected last 5 bytes %q, got: %q", expected, output)
	}
	if exitCode != 0 {
		t.Errorf("--bytes=5: Expected exit code 0, got %d", exitCode)
	}

	// Test with -c +N (start from byte N)
	output, exitCode, err = runTail(t, "-c", "+20", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-c +20: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "TUVWXYZ" // Bytes from 20th (index 19) onwards
	if output != expected {
		t.Errorf("-c +20: Expected bytes from 20 onwards %q, got: %q", expected, output)
	}
	if exitCode != 0 {
		t.Errorf("-c +20: Expected exit code 0, got %d", exitCode)
	}

	// Test with -c +N where N > byte count (should print nothing)
	output, exitCode, err = runTail(t, "-c", "+27", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-c +27: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = ""
	if output != expected {
		t.Errorf("-c +27: Expected empty output, got: %q", output)
	}
	if exitCode != 0 {
		t.Errorf("-c +27: Expected exit code 0, got %d", exitCode)
	}

	// Test combined number flag: -c15
	output, exitCode, err = runTail(t, "-c15", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-c15: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "LMNOPQRSTUVWXYZ" // Last 15 bytes
	if output != expected {
		t.Errorf("-c15: Expected last 15 bytes %q, got: %q", expected, output)
	}
	if exitCode != 0 {
		t.Errorf("-c15: Expected exit code 0, got %d", exitCode)
	}
}

func TestTailReverseOrder(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	lines := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	content := strings.Join(lines, "\n") + "\n"
	testFiles := []testFile{{name: "test.txt", content: []byte(content), mode: 0644}}
	createTestFiles(t, tmpDir, testFiles)

	// Test -r (reverse all lines by default)
	output, exitCode, err := runTail(t, "-r", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-r: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected := "Fifth\nFourth\nThird\nSecond\nFirst\n"
	if output != expected {
		t.Errorf("-r: Expected reversed lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-r: Expected exit code 0, got %d", exitCode)
	}

	// Test -r with -n N (reverse last N lines)
	output, exitCode, err = runTail(t, "-r", "-n", "3", "test.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("-r -n 3: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Fifth\nFourth\nThird\n" // Last 3 lines, reversed
	if output != expected {
		t.Errorf("-r -n 3: Expected last 3 reversed lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("-r -n 3: Expected exit code 0, got %d", exitCode)
	}
}

func TestTailWithMultipleFiles(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create test files
	testFiles := []testFile{
		{name: "file1.txt", content: []byte("A1\nB1\nC1\nD1\nE1\nF1\nG1\nH1\nI1\nJ1\nK1\n"), mode: 0644}, // 11 lines
		{name: "file2.txt", content: []byte("A2\nB2\nC2\n"), mode: 0644},                                 // 3 lines
	}
	createTestFiles(t, tmpDir, testFiles)

	// Test with multiple files, default tail (last 10 lines)
	output, exitCode, err := runTail(t, "file1.txt", "file2.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("Multi-file default: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}

	expectedFile1 := "B1\nC1\nD1\nE1\nF1\nG1\nH1\nI1\nJ1\nK1\n" // Last 10 lines of file1
	expectedFile2 := "A2\nB2\nC2\n"                             // All 3 lines of file2
	expected := fmt.Sprintf("==> %s <==\n%s\n==> %s <==\n%s", "file1.txt", expectedFile1, "file2.txt", expectedFile2)

	if output != expected {
		t.Errorf("Multi-file default: Expected output with headers, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Multi-file default: Expected exit code 0, got %d", exitCode)
	}

	// Test with multiple files and -n 2
	output, exitCode, err = runTail(t, "-n", "2", "file1.txt", "file2.txt")
	if err != nil && exitCode != 0 {
		t.Fatalf("Multi-file -n 2: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}

	expectedFile1 = "J1\nK1\n" // Last 2 lines of file1
	expectedFile2 = "B2\nC2\n" // Last 2 lines of file2
	expected = fmt.Sprintf("==> %s <==\n%s\n==> %s <==\n%s", "file1.txt", expectedFile1, "file2.txt", expectedFile2)

	if output != expected {
		t.Errorf("Multi-file -n 2: Expected output with headers and 2 lines each, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Multi-file -n 2: Expected exit code 0, got %d", exitCode)
	}
}

func TestTailWithNonExistentFile(t *testing.T) {
	// Run tail on a non-existent file
	output, exitCode, _ := runTail(t, "nonexistent.txt")

	// Check that an error message is printed to stderr (captured by CombinedOutput)
	if !strings.Contains(output, "tail: nonexistent.txt: No such file or directory") {
		t.Errorf("Expected error message about non-existent file, got: %s", output)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	// Test with one existing and one non-existing file
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()
	testFiles := []testFile{{name: "exists.txt", content: []byte("line1\nline2\n"), mode: 0644}}
	createTestFiles(t, tmpDir, testFiles)

	output, exitCode, _ = runTail(t, "exists.txt", "nonexistent.txt")

	// Check for header for existing file, its content, and the error message
	expectedHeader := "==> exists.txt <=="
	expectedContent := "line1\nline2\n"
	expectedError := "tail: nonexistent.txt: No such file or directory"

	if !strings.Contains(output, expectedHeader) {
		t.Errorf("Expected header for existing file, got:\n%s", output)
	}
	// Check content *before* the error message part
	if !strings.Contains(strings.Split(output, expectedError)[0], expectedContent) {
		t.Errorf("Expected content for existing file, got:\n%s", output)
	}
	if !strings.Contains(output, expectedError) {
		t.Errorf("Expected error message for non-existent file, got:\n%s", output)
	}

	// Exit code should still be 1 because one file failed
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 when one file is missing, got %d", exitCode)
	}
}

func TestTailWithPipedInput(t *testing.T) {
	stdinContent := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\nLine 12\n" // 12 lines

	// Default tail (last 10 lines)
	output, exitCode, err := runTailWithInput(t, stdinContent)
	if err != nil && exitCode != 0 {
		t.Fatalf("Stdin default: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected := "Line 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\nLine 12\n"
	if output != expected {
		t.Errorf("Stdin default: Expected last 10 lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Stdin default: Expected exit code 0, got %d", exitCode)
	}

	// Tail with -n 5
	output, exitCode, err = runTailWithInput(t, stdinContent, "-n", "5")
	if err != nil && exitCode != 0 {
		t.Fatalf("Stdin -n 5: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 8\nLine 9\nLine 10\nLine 11\nLine 12\n"
	if output != expected {
		t.Errorf("Stdin -n 5: Expected last 5 lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Stdin -n 5: Expected exit code 0, got %d", exitCode)
	}

	// Tail with -n +9 (start from line 9)
	output, exitCode, err = runTailWithInput(t, stdinContent, "-n", "+9")
	if err != nil && exitCode != 0 {
		t.Fatalf("Stdin -n +9: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 9\nLine 10\nLine 11\nLine 12\n"
	if output != expected {
		t.Errorf("Stdin -n +9: Expected lines from 9 onwards, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Stdin -n +9: Expected exit code 0, got %d", exitCode)
	}

	// Tail with -r (reverse)
	output, exitCode, err = runTailWithInput(t, stdinContent, "-r", "-n", "10")
	if err != nil && exitCode != 0 {
		t.Fatalf("Stdin -r: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 12\nLine 11\nLine 10\nLine 9\nLine 8\nLine 7\nLine 6\nLine 5\nLine 4\nLine 3\n"
	if output != expected {
		t.Errorf("Stdin -r: Expected reversed lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Stdin -r: Expected exit code 0, got %d", exitCode)
	}

	// Tail with -r -n 4 (reverse last 4)
	output, exitCode, err = runTailWithInput(t, stdinContent, "-r", "-n", "4")
	if err != nil && exitCode != 0 {
		t.Fatalf("Stdin -r -n 4: Expected no error, got exitCode %d, err: %v, output:\n%s", exitCode, err, output)
	}
	expected = "Line 12\nLine 11\nLine 10\nLine 9\n" // Last 4 lines, reversed
	if output != expected {
		t.Errorf("Stdin -r -n 4: Expected last 4 reversed lines, got:\n%s", output)
	}
	if exitCode != 0 {
		t.Errorf("Stdin -r -n 4: Expected exit code 0, got %d", exitCode)
	}
}

func TestTailFollow(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	fileName := "follow_test.txt"
	fileContent := "Initial Line 1\nInitial Line 2\n"
	testFiles := []testFile{{name: fileName, content: []byte(fileContent), mode: 0644}}
	createTestFiles(t, tmpDir, testFiles)
	filePath := filepath.Join(tmpDir, fileName)

	executable := getTailExecutable(t)
	cmd := exec.Command(executable, "-f", filePath)

	// Use pipes to read stdout incrementally
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start tail -f: %v", err)
	}

	// Use a channel to signal when reading is done or timed out
	outputChan := make(chan string)
	errChan := make(chan string)
	doneChan := make(chan bool)

	// Goroutine to read stdout
	go func() {
		buf := make([]byte, 1024)
		var output strings.Builder
		for {
			n, err := stdoutPipe.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
				// Send partial output or wait until done? Let's send full when done.
			}
			if err == io.EOF || err == io.ErrClosedPipe {
				break
			}
			if err != nil {
				t.Errorf("Error reading stdout pipe: %v", err)
				break
			}
		}
		outputChan <- output.String()
	}()

	// Goroutine to read stderr
	go func() {
		buf := make([]byte, 1024)
		var output strings.Builder
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err == io.EOF || err == io.ErrClosedPipe {
				break
			}
			if err != nil {
				t.Errorf("Error reading stderr pipe: %v", err)
				break
			}
		}
		errChan <- output.String()
	}()

	// Wait a moment for tail to start and print initial lines
	time.Sleep(100 * time.Millisecond) // Adjust timing as needed

	// Append data to the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("Failed to open file for appending: %v", err)
	}
	appendedText := "Appended Line 1\nAppended Line 2\n"
	if _, err := f.WriteString(appendedText); err != nil {
		f.Close()
		cmd.Process.Kill()
		t.Fatalf("Failed to write appended data: %v", err)
	}
	f.Close()

	// Wait a moment for tail to process the appended data
	time.Sleep(100 * time.Millisecond) // Adjust timing as needed

	// Stop the tail process
	if err := cmd.Process.Kill(); err != nil {
		// Ignore errors if the process already exited
		if !strings.Contains(err.Error(), "process already finished") {
			t.Logf("Error killing tail process: %v", err)
		}
	}

	// Wait for reading goroutines to finish
	cmd.Wait()      // Ensure the process resource is released
	close(doneChan) // Signal potential timeouts (though not strictly implemented here)

	stdoutOutput := <-outputChan
	stderrOutput := <-errChan

	if stderrOutput != "" {
		t.Errorf("Expected no stderr output, got: %s", stderrOutput)
	}

	// Expected output: Default is last 10 lines, so initially both initial lines.
	// Then the appended lines.
	expectedOutput := fileContent + appendedText
	if !strings.Contains(stdoutOutput, expectedOutput) { // Use contains as timing might affect exact match start
		t.Errorf("Expected stdout to contain:\n%s\nGot:\n%s", expectedOutput, stdoutOutput)
	}

	// Add more follow tests: e.g., with -n N, multiple appends, file rotation (harder)
}
