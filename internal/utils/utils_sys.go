package utils

import (
	"os"
	"strings"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"
const DASH_PATH = "internal/test_helpers/dash/your_shell.sh"

func IsTestingTesterUsingBusyboxOnAlpine(executablePath string) bool {
	isTestingTesterUsingBusybox := strings.HasSuffix(executablePath, ASH_PATH) || strings.HasSuffix(executablePath, DASH_PATH)

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusybox && isOnAlpine
}
