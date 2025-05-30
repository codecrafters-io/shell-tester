package internal

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testP3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "cat", CommandName: CUSTOM_CAT_COMMAND, CommandMetadata: ""},
		{CommandType: "head", CommandName: CUSTOM_HEAD_COMMAND, CommandMetadata: ""},
		{CommandType: "wc", CommandName: CUSTOM_WC_COMMAND, CommandMetadata: ""},
	}, false)
	if err != nil {
		return err
	}

	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Test-1
	randomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}
	filePath := path.Join(randomDir, fmt.Sprintf("file-%d", random.RandomInt(1, 100)))
	randomWords := random.RandomWords(5)
	fileContent := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", randomWords[0], randomWords[1], randomWords[2], randomWords[3], randomWords[4])
	if err := writeFiles([]string{filePath}, []string{fileContent}, logger); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	lines := strings.Count(fileContent, "\n") + 1
	words := strings.Count(strings.ReplaceAll(fileContent, "\n", " "), " ") + 1
	bytes := len(fileContent)

	input := fmt.Sprintf(`cat %s | head -n 5 | wc`, filePath)
	expectedOutput := fmt.Sprintf("%8d%8d%8d", lines, words, bytes)

	singleLineTestCase := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: nil,
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Test-2
	newRandomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return err
	}
	randomUniqueFileNames := random.RandomInts(1, 100, 6)
	filePaths := []string{
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[0])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[1])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[2])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[3])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[4])),
		path.Join(newRandomDir, fmt.Sprintf("f-%d", randomUniqueFileNames[5])),
	}
	fileContents := random.RandomWords(6)
	if err := writeFiles(filePaths, fileContents, logger); err != nil {
		return err
	}

	sort.Ints(randomUniqueFileNames)
	availableEntries := randomUniqueFileNames[1:4]

	input = fmt.Sprintf(`ls -la %s | tail -n 5 | head -n 3 | grep "f-%d"`, newRandomDir, availableEntries[2])
	expectedRegexPattern := fmt.Sprintf("^[rwx-]* .* f-%d", availableEntries[1])

	singleLineTestCase2 := test_cases.CommandResponseTestCase{
		Command:          input,
		ExpectedOutput:   "", // We completely rely on the regex pattern here
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(expectedRegexPattern)},
		SuccessMessage:   "✓ Received expected output",
	}
	if err := singleLineTestCase2.Run(asserter, shell, logger); err != nil {
		fileName := path.Join(newRandomDir, fmt.Sprintf("f-%d", availableEntries[2]))
		// expectedOutput := "-rw-r--r-- 1 root root    9 May 30 13:21 f-32"
		fmt.Println(generateLSOutput(fileName))

		return err
	}

	return logAndQuit(asserter, nil)
}

func formatPermissions(mode os.FileMode) string {
	var result string

	// File type
	switch {
	case mode.IsDir():
		result += "d"
	case mode&os.ModeSymlink != 0:
		result += "l"
	case mode&os.ModeDevice != 0:
		result += "b"
	case mode&os.ModeCharDevice != 0:
		result += "c"
	case mode&os.ModeNamedPipe != 0:
		result += "p"
	case mode&os.ModeSocket != 0:
		result += "s"
	default:
		result += "-"
	}

	// Owner permissions
	perm := mode.Perm()
	result += formatTriad(perm, 6) // bits 8,7,6
	result += formatTriad(perm, 3) // bits 5,4,3
	result += formatTriad(perm, 0) // bits 2,1,0

	return result
}

func formatTriad(perm os.FileMode, shift uint) string {
	var result string

	// Read permission
	if perm&(4<<shift) != 0 {
		result += "r"
	} else {
		result += "-"
	}

	// Write permission
	if perm&(2<<shift) != 0 {
		result += "w"
	} else {
		result += "-"
	}

	// Execute permission
	if perm&(1<<shift) != 0 {
		result += "x"
	} else {
		result += "-"
	}

	return result
}

func getUsername(uid uint32) string {
	u, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return strconv.Itoa(int(uid))
	}
	return u.Username
}

func getGroupname(gid uint32) string {
	g, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return strconv.Itoa(int(gid))
	}
	return g.Name
}

func generateLSOutput(fileName string) string {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}

	// Get system-specific info
	sys := fileInfo.Sys().(*syscall.Stat_t)

	// Format permissions
	mode := fileInfo.Mode()
	permissions := formatPermissions(mode)

	// Get owner and group names
	ownerName := getUsername(sys.Uid)
	groupName := getGroupname(sys.Gid)

	// Format modification time
	modTime := fileInfo.ModTime().Format("Jan 2 15:04")

	// Print in ls -la format
	return fmt.Sprintf("%s %2d %s %s %8d %s %s\n",
		permissions,
		sys.Nlink, // number of hard links
		ownerName,
		groupName,
		fileInfo.Size(),
		modTime,
		fileInfo.Name(),
	)
}
