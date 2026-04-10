package utils

import (
	"os"
	"strings"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"
const DASH_PATH = "internal/test_helpers/dash/your_shell.sh"

func IsTestingTesterUsingBusyboxOnAlpine(shellPath string) bool {
	isTestingTesterUsingBusybox :=
		strings.HasSuffix(shellPath, ASH_PATH) || strings.HasSuffix(shellPath, DASH_PATH)

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusybox && isOnAlpine
}
