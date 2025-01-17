package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/tester-utils/random"
)

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

const CUSTOM_LS_COMMAND = "ls"
const CUSTOM_CAT_COMMAND = "cat"

func getRandomInvalidCommand() string {
	return getRandomInvalidCommands(1)[0]
}

func getRandomInvalidCommands(n int) []string {
	words := random.RandomWords(n)
	invalidCommands := make([]string, n)

	for i := 0; i < n; i++ {
		invalidCommands[i] = "invalid_" + words[i] + "_command"
	}

	return invalidCommands
}

func getRandomString() string {
	// We will use a random numeric string of length = 6
	var result string
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func getRandomName() string {
	names := []string{"Alice", "David", "Emily", "James", "Maria"}
	return names[random.RandomInt(0, len(names))]
}

func logAndQuit(asserter *logged_shell_asserter.LoggedShellAsserter, err error) error {
	asserter.LogRemainingOutput()
	return err
}
