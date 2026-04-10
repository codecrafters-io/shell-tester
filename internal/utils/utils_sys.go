package utils

import (
	"os"
	"strings"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"
const DASH_PATH = "internal/test_helpers/dash/your_shell.sh"
const BASH_PATH = "internal/test_helpers/bash/your_shell.sh"

func IsTestingTester(executablePath string) bool {
	return (strings.HasSuffix(executablePath, ASH_PATH) ||
		strings.HasSuffix(executablePath, DASH_PATH) ||
		strings.HasSuffix(executablePath, BASH_PATH))
}

func IsTestingTesterUsingBusyboxOnAlpine(executablePath string) bool {
	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return IsTestingTester(executablePath) && isOnAlpine
}
